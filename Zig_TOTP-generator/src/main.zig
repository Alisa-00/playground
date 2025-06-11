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

    //define a dir/filename to read from. maybe just .vault.json ?
    const vault_directory = std.fs.cwd();
    const vault_filename = "vault.json";

    var args = std.process.args();
    const cmd = args.next().?;

    // Maybe try a different allocator? arena?
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();

    // Create vault file if it doesnt exist
    try vault.checkVault(vault_directory, vault_filename);

    const action_arg = args.next() orelse exitError(Error.MissingArguments, cmd, true);
    const action = readAction(action_arg);
    const name_arg: ?[]const u8 = args.next();

    // Read any entries from vault file
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
            //handle the null case properly, print appropriate msg for each case
            const secret = vlt.get(name).?;

            //const totp: []const u8 = getTotp(name, secret) catch |err| exitError(err);
            try stdout.print("{s}:{s}\n", .{name, secret});
        },
        Action.List => {
            var iter = vlt.keyIterator();
            while (iter.next()) |key| {
                // get totp before displaying
                const value = vlt.get(key.*).?;
                try stdout.print("Name: {s}\nSecret: {s}\n\n", .{key.*, value});
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
