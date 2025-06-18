const std = @import("std");
const hmac = std.crypto.auth.hmac;

pub fn generateTOTP(secret: []const u8, period: i64, digits: u32) u32 {
    const unix_time = std.time.timestamp();
    const sequence_value = @divFloor(unix_time, period);

    var msg: [8]u8 = undefined;
    std.mem.writeInt(u64, &msg, @intCast(sequence_value), std.builtin.Endian.big);

    var hmac_sha1 = hmac.HmacSha1.init(secret);
    hmac_sha1.update(&msg);

    var result: [hmac.HmacSha1.mac_length]u8 = undefined;
    hmac_sha1.final(&result);

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

