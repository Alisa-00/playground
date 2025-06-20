const std = @import("std");
const hmac = std.crypto.auth.hmac;

const BASE32 = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567";

pub const decode_error = error {
    NotBase32Character,
};

fn upper(text: []const u8, allocator: std.mem.Allocator) ![]u8 {
    var buffer: []u8 = try allocator.alloc(u8, text.len);
    for (0..text.len) |i| {
        var uppercase_char: u8 = text[i];

        if (text[i] >= 'a' and text[i] <= 'z') {
            uppercase_char -= 32;
        }

        buffer[i] = uppercase_char;
    }
    return buffer;
}

fn base32Index(char: u8) !u8 {
    if (char == '=') return 0;

    for (0..31) |i| {
        if (char == BASE32[i]) return @intCast(i);
    }

    return decode_error.NotBase32Character;
}

fn decode(secret: []const u8, allocator: std.mem.Allocator) ![]u8 {
    const uppercase_secret = try upper(secret, allocator);
    defer allocator.free(uppercase_secret);

    const decoded_len = (secret.len*5)/8;
    var buffer: []u8 = try allocator.alloc(u8, decoded_len);
    @memset(buffer, 32);

    var in_idx: usize = 0;
    var out_idx: usize = 0;

    while (in_idx < uppercase_secret.len) : (in_idx += 8) {

        const in_byte1 = try base32Index(uppercase_secret[in_idx]);
        const in_byte2 = try base32Index(uppercase_secret[in_idx+1]);
        const in_byte3 = try base32Index(uppercase_secret[in_idx+2]);
        const in_byte4 = try base32Index(uppercase_secret[in_idx+3]);
        const in_byte5 = try base32Index(uppercase_secret[in_idx+4]);
        const in_byte6 = try base32Index(uppercase_secret[in_idx+5]);
        const in_byte7 = try base32Index(uppercase_secret[in_idx+6]);
        const in_byte8 = try base32Index(uppercase_secret[in_idx+7]);

        const out_byte1 = (in_byte1 << 3) | (in_byte2 >> 2);
        const out_byte2 = (in_byte2 << 6) | (in_byte3 << 1) | (in_byte4 >> 4);
        const out_byte3 = (in_byte4 << 4) | (in_byte5 >> 1);
        const out_byte4 = (in_byte5 << 7) | (in_byte6 << 2) | (in_byte7 >> 3);
        const out_byte5 = (in_byte7 << 5) | (in_byte8);

        buffer[out_idx]   = out_byte1;
        buffer[out_idx+1] = out_byte2;
        buffer[out_idx+2] = out_byte3;
        buffer[out_idx+3] = out_byte4;
        buffer[out_idx+4] = out_byte5;

        out_idx += 5;
    }

    return buffer;

}

pub fn generateTOTP(secret: []const u8, period: i64, digits: u32) !u32 {
    var gpa = std.heap.GeneralPurposeAllocator(.{}){}; 
    const allocator = gpa.allocator();

    const unix_time = std.time.timestamp();
    const sequence_value = @divFloor(unix_time, period);

    var msg: [8]u8 = undefined;
    std.mem.writeInt(u64, &msg, @intCast(sequence_value), std.builtin.Endian.big);

    const decoded_secret = try decode(secret, allocator);
    defer allocator.free(decoded_secret);

    var hmac_sha1 = hmac.HmacSha1.init(decoded_secret);
    hmac_sha1.update(&msg);

    var result: [hmac.HmacSha1.mac_length]u8 = undefined;
    hmac_sha1.final(&result);

    const truncated_result = dynamicTruncate(&result);
    const output_digits = truncated_result % std.math.pow(u32, 10, digits);

    return output_digits;
}

fn dynamicTruncate(string: []const u8) u32 {
    const offset = string[string.len-1] & 0x0F;

    const byte1 = @as(u32, string[offset]) << 24;
    const byte2 = @as(u32, string[offset+1]) << 16;
    const byte3 = @as(u32, string[offset+2]) << 8;
    const byte4 = @as(u32, string[offset+3]);

    const bytes = byte1 | byte2 | byte3 | byte4;
    const result_masked = bytes & 0x7FFFFFFF;
    return result_masked;
}
