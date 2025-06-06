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
    _ = args.next().?;

    // Maybe try a different allocator? arena?
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();

    // Handle any errors, but mostly create the file if it doesnt exist
    var vlt = std.StringHashMap([]const u8).init(allocator);
    try vault.readHMap(vault_directory, vault_filename, &vlt, allocator);
    // no need to defer free, program execution ends when this scope exits

    const action_arg = args.next() orelse exitError(Error.MissingArguments);
    const action = readAction(action_arg);

    switch(action) {
        Action.Add => {
            const name = args.next() orelse exitError(Error.MissingArguments);
            const secret = args.next() orelse exitError(Error.MissingArguments);

            if (!vlt.contains(name)) {
                try vlt.put(name, secret);
                try vault.writeHMap(vault_directory, vault_filename, vlt, allocator);
                try stdout.print("Added {s} to vault successfully!\n", .{name});
            } else {
                try stdout.print("Could not add '{s}'. already exists in vault!\n", .{name});
            }
        },
        Action.Get => {
            const name = args.next() orelse exitError(Error.MissingArguments);
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
        Action.Unknown => {},
    }



}

fn readAction(arg: []const u8) Action {

    if (std.mem.eql(u8, arg, "add")) { return Action.Add; }
    else if (std.mem.eql(u8, arg, "get")) { return Action.Get; }
    else if (std.mem.eql(u8, arg, "list")) { return Action.List; }
    return Action.Unknown;

}

fn exitError(err: anyerror) noreturn {

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
    std.process.exit(0);
}

test "simple test" {
    var list = std.ArrayList(i32).init(std.testing.allocator);
    defer list.deinit(); // Try commenting this out and see if zig detects the memory leak!
    try list.append(42);
    try std.testing.expectEqual(@as(i32, 42), list.pop());
}
