const std = @import("std");

pub fn checkVault(dir: std.fs.Dir, filename: []const u8) !void {
    dir.access(filename, .{}) catch |err| {
        if (err == std.posix.AccessError.FileNotFound) {

            std.debug.print("Vault file not found. creating new vault file...\n", .{});

            const file = try dir.createFile(filename, .{});
            const written = try file.write("{}");

            if (written > 0) {
                std.debug.print("New empty vault file created!\n", .{});
            }

            file.close();
        }
    };
}

pub fn readHMap(dir: std.fs.Dir, filename: []const u8, hmap: *std.StringHashMap([]const u8), allocator: std.mem.Allocator) !void {

    const json_file = try dir.openFile(filename, .{});
    defer json_file.close();

    const file_metadata = try json_file.metadata();
    const file_size = file_metadata.size();
    const buffer = try json_file.readToEndAlloc(allocator, file_size);
    defer allocator.free(buffer);

    const parse = try std.json.parseFromSlice(std.json.Value, allocator, buffer, .{});
    defer parse.deinit();

    const obj = parse.value;

    var iter = obj.object.iterator();
    while (iter.next()) |entry| {
        const key = try allocator.dupe(u8, entry.key_ptr.*);
        const val = try allocator.dupe(u8, entry.value_ptr.*.string);

        try hmap.*.put(key, val);
    }
}

test "read HMap from file" {

    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();

    const dir = try std.fs.cwd().makeOpenPath("test", .{});
    const filename: []const u8 = "test.json";

    var hmap = std.StringHashMap([]const u8).init(allocator);
    try readHMap(dir, filename, &hmap, allocator);

    var iter = hmap.keyIterator();
    while (iter.next()) |key| {
        const value = hmap.get(key.*).?;
        std.debug.print("{s}:{s}\n", .{key.*, value});
    } 

}

pub fn writeHMap(dir: std.fs.Dir, filename: []const u8, hmap: std.StringHashMap([]const u8), allocator: std.mem.Allocator) !void {

    var buffer = std.ArrayList(u8).init(allocator);
    defer buffer.deinit();

    var json_file = try dir.createFile(filename, .{});
    defer json_file.close();

    const writer = json_file.writer();
    var json_stream = std.json.writeStream(writer, .{});
    try json_stream.beginObject();

    var it = hmap.iterator();
    while (it.next()) |entry| {
        try json_stream.objectField(entry.key_ptr.*);
        try json_stream.write(entry.value_ptr.*);
    }

    try json_stream.endObject();
}

pub fn cleanVault(vlt: std.StringHashMap([]const u8), name: ?[]const u8, allocator: std.mem.Allocator) void {
    var iter = vlt.iterator();
    const null_name = name == null;

    while (iter.next()) |acc| {
        const key = @constCast(acc.key_ptr.*);
        const val = @constCast(acc.value_ptr.*);

        @memset(key, 0);
        @memset(val, 0);

        if (!null_name) {
            if (std.mem.eql(u8, name.?, key)) continue;
        }
        allocator.free(key);
        allocator.free(val);
    }
}

test "write test" {

    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();

    const dir = try std.fs.cwd().makeOpenPath("test", .{});
    const filename: []const u8 = "test.json";

    const acc1 = "cuenta1";
    const acc2 = "cuenta2";
    const acc3 = "cuenta3";
    const sec1 = "CUENTA MIA";
    const sec2 = "CUENTA NIA";
    const sec3 = "CUENTA NYA";

    var hmap = std.StringHashMap([]const u8).init(allocator);

    try hmap.put(acc1, sec1);
    try hmap.put(acc2, sec2);
    try hmap.put(acc3, sec3);

    try writeHMap(dir, filename, hmap, allocator);

    hmap.clearAndFree();
    try readHMap(dir, filename, &hmap, allocator);

    var iter = hmap.keyIterator();
    while (iter.next()) |key| {
        const value = hmap.get(key.*).?;
        std.debug.print("{s}:{s}\n", .{key.*, value});
    }
}
