const std = @import("std");

const EncodeDecodeError = error { 
    EncodeError,
    DecodeError,
    UnknownAction,
    NotBase64Character,
    Other
};

const Action = enum {
    Encode,
    Decode,
    Unknown
};

const ERROR_ENCODE = "Error: encoding failed";
const ERROR_DECODE = "Error: decoding failed";
const ERROR_ACTION = "Error: unknown action";
const ERROR_OTHERS = "Unknown error ocurred";
const ERROR_BASE64 = "Error: Input is not a base64 encoded string";

const BASE64_TABLE = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

pub fn main() !void {

    var gpa = std.heap.GeneralPurposeAllocator(.{}){};
    const allocator = gpa.allocator();

    // read and validate cli arguments
    var args = std.process.args();
    const cmd = args.next().?;

    const action = args.next();
    if (action == null) {
        printUsage(cmd);
        return;
    }

    const data = args.next();
    if (data == null) {
        printUsage(cmd);
        return;
    }

    const action_value: Action = readAction(action.?);

    switch (action_value) {
        Action.Encode => {
            b64Encode(data.?, allocator) catch |err| handleError(err);
        },
        Action.Decode => {
            b64Decode(data.?, allocator) catch |err| handleError(err);
        },
        Action.Unknown => {
            handleError(EncodeDecodeError.UnknownAction);
        }
    }

}

fn readAction (action: []const u8) Action {

    if (std.mem.eql(u8, action, "encode")) { return Action.Encode; }
    else if (std.mem.eql(u8, action, "decode")) { return Action.Decode; }

    return Action.Unknown;

}

fn b64Encode (text: []const u8, allocator: std.mem.Allocator) !void {

    // calculate padding bytes to reach length divisible by 6
    var modulo = text.len*8;
    modulo = modulo % 6;

    const pad_bytes: u8 = switch(modulo) {
        2 => 2, 4 => 1, 0 => 0,
        else => unreachable
    };

    // copy our data into padded array, pad bytes stay as 0
    var padded_array = try allocator.alloc(u8, text.len+pad_bytes);
    defer allocator.free(padded_array);
    @memset(padded_array, 0);
    @memcpy(padded_array[0..text.len], text);

    // Initialize string to hold the base 64 representation
    const b64_len = padded_array.len + (padded_array.len/3);
    var b64_string = try allocator.alloc(u8, b64_len);
    defer allocator.free(b64_string);
    @memset(b64_string, 0);

    var idx: u32 = 0;
    var out_idx: u32 = 0;
    // padded array.len is always a multiple of 3 so this is safe
    while (idx < padded_array.len) : (idx += 3) {

        // 0bXXXXXXYY 0bYYYYZZZZ 0bZZAAAAAA
        const byte1: u8 = padded_array[idx];
        const byte2: u8 = padded_array[idx+1];
        const byte3: u8 = padded_array[idx+2];

        // keep 6 XXXXXX
        const outbyte1: u8 = byte1 >> 2;

        // keep 2 YY from byte1, YYYY from byte 2. shift them into position and sum
        var outbyte2: u8 = byte1 & 0b00000011;
        outbyte2 = outbyte2 << 4;
        outbyte2 = outbyte2 + ((byte2 & 0b11110000) >> 4);

        // similar to outbyte2, ZZZZ from byte2, ZZ from byte3
        var outbyte3: u8 = byte2 & 0b00001111;
        outbyte3 = outbyte3 << 2;
        outbyte3 =  outbyte3 + ((byte3 & 0b11000000) >> 6);

        // siilar to outbyte1 but no shift needed
        const outbyte4: u8 = byte3 & 0b00111111;

        // write into b64_string
        b64_string[out_idx] = BASE64_TABLE[outbyte1];
        b64_string[out_idx+1] = BASE64_TABLE[outbyte2];
        b64_string[out_idx+2] = BASE64_TABLE[outbyte3];
        b64_string[out_idx+3] = BASE64_TABLE[outbyte4];

        out_idx += 4;
    }

    // overwrite padding bytes with =
    if (pad_bytes > 0) {
        const unpadded_len = b64_len - pad_bytes;
        @memset(b64_string[unpadded_len..], '=');
    }

    std.debug.print("{s}\n", .{b64_string});
}

fn b64Decode(text: []const u8, allocator: std.mem.Allocator) !void {

    // 3 chars per 4 b64 chars
    var out_len = 3 * text.len / 4;

    if (text[text.len-2] == '=') { out_len -=2; }
    else if (text[text.len-1] == '=') { out_len -= 1; }

    var out_array = try allocator.alloc(u8, out_len);
    defer allocator.free(out_array);
    @memset(out_array, 0);


    var idx: u8 = 0;
    var out_idx: u8 = 0;
    while (idx < text.len) : (idx += 4) {

        const aux2 = if (idx+2 >= text.len) '=' else text[idx+2];
        const aux3 = if (idx+3 >= text.len) '=' else text[idx+3];

        const byte1: u8 = base64Index(text[idx]) catch |err| handleError(err);
        const byte2: u8 = base64Index(text[idx+1]) catch |err| handleError(err);
        const byte3: u8 = base64Index(aux2) catch |err| handleError(err);
        const byte4: u8 = base64Index(aux3) catch |err| handleError(err);

        var outbyte1: u8 = byte1 << 2;
        outbyte1 = outbyte1 + (byte2 >> 4);

        var outbyte2: u8 = byte2 << 4;
        outbyte2 = outbyte2 + (byte3 >> 2);

        var outbyte3: u8 = byte3 << 6;
        outbyte3 = outbyte3 + byte4;

        out_array[out_idx] = outbyte1;

        if (out_idx+1 >= out_array.len) break;
        out_array[out_idx+1] = outbyte2;

        if (out_idx+2 >= out_array.len) break;
        out_array[out_idx+2] = outbyte3;

        out_idx += 3;
    }

    std.debug.print("{s}\n", .{out_array});
}

fn base64Index(char: u8) !u8 {

    if (char == '=') return 64;

    for (0..63) |i| {
        if (char == BASE64_TABLE[i]) return @intCast(i);
    }

    return EncodeDecodeError.NotBase64Character;
}

fn handleError(err: anyerror) noreturn {

    var msg: []const u8 = undefined;

    switch(err) {
        EncodeDecodeError.EncodeError => { msg = ERROR_ENCODE; },
        EncodeDecodeError.DecodeError => { msg = ERROR_DECODE; },
        EncodeDecodeError.UnknownAction => { msg = ERROR_ACTION; },
        EncodeDecodeError.NotBase64Character => { msg = ERROR_BASE64; },
        EncodeDecodeError.Other => { msg = ERROR_OTHERS; },
        else => { msg = ERROR_OTHERS; },
    }

    std.debug.print("{s}\n", .{msg});
    std.process.exit(1);
}

fn printUsage(cmd: []const u8) void {
    std.debug.print("Usage: {s} [encode/decode] [data]\n", .{cmd});
}
