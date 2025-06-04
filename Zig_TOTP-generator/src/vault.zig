const std = @import("std");

// Vault data structure should be a Hash map.
// Each entry is name:secret
// Json saved as name:secret as well
pub const Vault = struct {
    accounts: []Account,
    parse: ?std.json.Parsed([]Account),
    buffer: ?[]u8,
};

 // Unneeded - each acc is an entry in the vault
 // name and secret are already stored as key and value
pub const Account = struct {
    name: []const u8,
    secret: []const u8,
};

 // Clean depends on the hash map structure
 // No need for parse and buffer bs
pub fn cleanVault(vault: *const Vault, allocator: std.mem.Allocator) void {
    vault.parse.?.deinit();
    allocator.free(vault.buffer.?);
}

 // Storing accs into the hash map will likely need to be done 1 by 1
 // Can read file in 1 go, but still need to process accs 1 by 1
 // Research!!
pub fn readVault(dir: std.fs.Dir, filename: []const u8, allocator: std.mem.Allocator) !Vault {

    const json_file = try dir.openFile(filename, .{});
    defer json_file.close();

    const buffer = try json_file.readToEndAlloc(allocator, 10 * 1024); // 10 KB max
    const parse = try std.json.parseFromSlice([]Account, allocator, buffer, .{});

    return Vault {
        .accounts = parse.value,
        .parse = parse,
        .buffer = buffer,
    };

}

 // write name:secret, one per row ideally
pub fn writeVault(dir: std.fs.Dir, filename: []const u8, allocator: std.mem.Allocator, accs: []Account) !void {

    var buffer = std.ArrayList(u8).init(allocator);
    defer buffer.deinit();

    var json_file = try dir.createFile(filename, .{});
    defer json_file.close();

    try std.json.stringify(accs, .{}, buffer.writer());
    try json_file.writeAll(buffer.items);

}

test "write read test" {

    //const page_allocator = std.heap.page_allocator;
    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const gpa_allocator = gpa.allocator();
    const expect = std.testing.expect;
    const eql = std.mem.eql;

    const dir = try std.fs.cwd().makeOpenPath("test", .{});
    const filename: []const u8 = "test.json";

    const acc1 = Account { .name = "cuenta1", .secret = "CUENTA MIA" };
    const acc2 = Account { .name = "account", .secret = "quack" };
    const acc3 = Account { .name = "kurwa", .secret = "miku polish cow" };
    const acc4 = Account { .name = "akhi", .secret = "falafel hummus pita" };
    const accs = [_]Account{acc1, acc2, acc3, acc4};
    const ptrAccs = @constCast(accs[0..accs.len]);

    try writeVault(dir, filename, gpa_allocator, ptrAccs);

    const vault: Vault = try readVault(dir, filename, gpa_allocator);
    defer cleanVault(&vault, gpa_allocator);

    for (0..vault.accounts.len) |i| {
        try expect(eql(u8, accs[i].name, vault.accounts[i].name));
        try expect(eql(u8, accs[i].secret, vault.accounts[i].secret));
    }
}
