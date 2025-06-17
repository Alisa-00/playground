const std = @import("std");
const vault = @import("vault");
const stdout = std.io.getStdOut().writer();

const Action = enum {
    Add,
    Get,
    List,
    Unknown,
};

const Error = error {
    UnknownAction,
    MissingArguments,
    Other,
};


pub fn main() !void {

    const vault_directory = std.fs.cwd();
    const vault_filename = ".vaultfile";

    var args = std.process.args();
    const cmd = args.next().?;

    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();

    // Create vault file if it doesnt exist
    try vault.checkVault(vault_directory, vault_filename);

    const action_arg = args.next() orelse exitError(Error.MissingArguments, cmd, true);
    const action = readAction(action_arg);
    const name_arg: ?[]const u8 = args.next();

    // Read entries from vault file
    var vlt = std.StringHashMap([]const u8).init(allocator);
    try vault.readHMap(vault_directory, vault_filename, &vlt, allocator);
    defer vlt.deinit();
    defer vault.cleanVault(vlt, name_arg, allocator);


    switch(action) {
        Action.Add => {
            const name = name_arg orelse exitError(Error.MissingArguments, cmd, true);
            const secret = args.next() orelse exitError(Error.MissingArguments, cmd, true);

            if (!vlt.contains(name)) {
                try vlt.put(name, secret);
                try vault.writeHMap(vault_directory, vault_filename, vlt, allocator);
                try stdout.print("Added {s} to vault successfully!\n", .{name});
            } else {
                try stdout.print("Could not add '{s}'. already exists in vault!\n", .{name});
            }
        },
        Action.Get => {
            const name = name_arg orelse exitError(Error.MissingArguments, cmd, true);
            const secret = vlt.get(name).?;

            const period: i64 = 30;
            const digits: u32 = 6;
            const totp = generateTOTP(secret, period, digits);
            try stdout.print("{s} - {d}\n", .{name, totp});
        },
        Action.List => {
            var iter = vlt.keyIterator();
            while (iter.next()) |key| {
                const value = vlt.get(key.*).?;
                const period: i64 = 30;
                const digits: u32 = 6;
                const totp = generateTOTP(value, period, digits);
                try stdout.print("{s}: {d}\n", .{key.*, totp});
            }
        },
        Action.Unknown => {
            exitError(Error.UnknownAction, cmd, true);
        },
    }

}

fn readAction(arg: []const u8) Action {

    if (std.mem.eql(u8, arg, "add")) { return Action.Add; }
    else if (std.mem.eql(u8, arg, "get")) { return Action.Get; }
    else if (std.mem.eql(u8, arg, "list")) { return Action.List; }
    return Action.Unknown;

}

fn generateTOTP(secret: []const u8, period: i64, digits: u32) u32 {
    const unix_time = std.time.timestamp();
    const sequence_value = @divFloor(unix_time, period);

    var msg: [8]u8 = undefined;
    std.mem.writeInt(u64, &msg, @intCast(sequence_value), std.builtin.Endian.big);

    var hmac = std.crypto.auth.hmac.HmacSha1.init(secret);
    hmac.update(&msg);

    var result: [std.crypto.auth.hmac.HmacSha1.mac_length]u8 = undefined;
    hmac.final(&result);

    // dynamic truncation
    const offset = result[result.len-1] & 0x0F;

    const byte1 = @as(u32, result[offset]) << 24;
    const byte2 = @as(u32, result[offset+1]) << 16;
    const byte3 = @as(u32, result[offset+2]) << 8;
    const byte4 = @as(u32, result[offset+3]);

    const bytes = byte1 | byte2 | byte3 | byte4;
    const result_masked = bytes & 0x7FFFFFFF;

    const output_digits = result_masked % std.math.pow(u32, 10, digits);

    return output_digits;
}

fn exitError(err: anyerror, cmd: []const u8, print_usage: bool) noreturn {

    const ERROR_ACTION = "ERROR: Unknown action";
    const ERROR_ARGS = "ERROR: Missing arguments";
    const ERROR_OTHER = "ERROR: Something went wrong";

    const msg = switch(err) {
        Error.UnknownAction => ERROR_ACTION,
        Error.MissingArguments => ERROR_ARGS,
        Error.Other => ERROR_OTHER,
        else => ERROR_OTHER
    };

    std.debug.print("{s}\n", .{ msg });

    if (print_usage) {
        std.debug.print("Usage:\n{s} list\n", .{cmd});
        std.debug.print("{s} add [account]\n", .{cmd});
        std.debug.print("{s} get [account]\n", .{cmd});

    }

    std.process.exit(0);
}
