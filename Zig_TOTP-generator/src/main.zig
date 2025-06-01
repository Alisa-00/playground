const std = @import("std");
const file = @import("vault");
const stdout = std.io.getStdOut().writer();
const Account = file.Account;
const Vault = file.Vault;

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

const VAULT_DIR = std.fs.cwd();
const VAULT_FILENAME = "vault.json";

pub fn main() !void {

    var args = std.process.args();
    _ = args.next().?;

    const allocator = std.heap.page_allocator;

    const vlt: Vault = try file.readVault(VAULT_DIR, VAULT_FILENAME, allocator);
    defer file.cleanVault(vlt, allocator);


    const action_arg = args.next() orelse exitError(Error.MissingArguments);
    const action = readAction(action_arg);

    switch(action) {
        Action.Add => {
            const name = args.next() orelse exitError(Error.MissingArguments);
            const secret = args.next() orelse exitError(Error.MissingArguments);

            addAccount(name, secret, vlt, allocator) catch |err| exitError(err);
            try stdout.print("Added {s} to vault successfully!\n", .{name});
        },
        Action.Get => {
            const name = args.next() orelse exitError(Error.MissingArguments);

            //const totp: []const u8 = getTotp(name, vlt) catch |err| exitError(err);
            try stdout.print("{s}:\n", .{name});
        },
        Action.List => {
            for (vlt.accounts) |acc| {
                try stdout.print("Name: {s}\nSecret: {s}\n\n", .{acc.name, acc.secret});
            }
        },
        Action.Unknown => {},
    }


    try file.writeVault(VAULT_DIR, VAULT_FILENAME, allocator, vlt.accounts);

}

fn readAction(arg: []const u8) Action {

    if (std.mem.eql(u8, arg, "add")) { return Action.Add; }
    else if (std.mem.eql(u8, arg, "get")) { return Action.Get; }
    else if (std.mem.eql(u8, arg, "list")) { return Action.List; }
    return Action.Unknown;

}

fn addAccount(name: []const u8, secret: []const u8, vlt: Vault, allocator: std.mem.Allocator) !void {

    var new_accs = try allocator.alloc(Account, vlt.accounts.len+1);
    //comptime { std.debug.print("new_accs {any}\nvlt.accounts {any}\nvlt {any}", .{@TypeOf(new_accs), @TypeOf(vlt.accounts), @TypeOf(vlt)});}
    @memcpy(new_accs[0..vlt.accounts.len], vlt.accounts);
    new_accs[vlt.accounts.len] = Account { .name = name, .secret = secret };

    vlt.accounts = new_accs;

}

fn getVault(name: []const u8) Error!Vault {
    std.debug.print("{s}\n", .{ name });
    return Error.Other;
}

fn getVaults() Error!*[]Vault {
    return Error.Other;
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
