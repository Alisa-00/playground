const std = @import("std");

pub fn build(b: *std.Build) !void {

    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    const server_exe = b.addExecutable(.{
        .name = "dbserver",
        //.root_source_file = b.path("src/database/main.c"),
        .target = target,
        .optimize = optimize
    });

    server_exe.linkLibC();
    server_exe.root_module.addIncludePath(b.path("include"));
    server_exe.root_module.addIncludePath(b.path("../../../../../usr/include"));

    server_exe.addCSourceFiles(.{
        .files = &.{
            "src/database/main.c",
            "src/database/db_poll.c",
            "src/database/file.c",
            "src/database/parse.c",
        },
        .flags = &.{},
    });

    b.installArtifact(server_exe);

    const client_exe = b.addExecutable(.{
        .name = "dbclient",
        //.root_source_file = b.path("src/client/client.c"),
        .target = target,
        .optimize = optimize
    });

    client_exe.linkLibC();
    client_exe.root_module.addIncludePath(b.path("include"));
    client_exe.root_module.addIncludePath(b.path("../../../../../usr/include"));

    client_exe.addCSourceFiles(.{
        .files = &.{
            "src/client/client.c",
        },
        .flags = &.{},
    });

    b.installArtifact(client_exe);
}
