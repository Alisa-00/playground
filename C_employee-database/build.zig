const std = @import("std");

pub fn build(b: *std.Build) !void {

    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    const server_exe = b.addExecutable(.{
        .name = "dbserver",
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

    const run_server = b.addRunArtifact(server_exe);
    run_server.step.dependOn(b.getInstallStep());
    run_server.addArgs(
        &[_][]const u8{
            "-f",
            "./mytestdb.db",
            "-p",
            "5555",
            "-l",
        }
    );

    const server_run_step = b.step("run", "Run the server");
    server_run_step.dependOn(&run_server.step);

    const client_exe = b.addExecutable(.{
        .name = "dbclient",
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

    const run_client = b.addRunArtifact(client_exe);
    run_client.step.dependOn(b.getInstallStep());

    run_client.addArgs(
        &[_][]const u8{
            "-h",
            "127.0.0.1",
            "-p",
            "5555",
            "-l",
        }
    );

    const client_run_step = b.step("runclient", "Run the client");
    client_run_step.dependOn(&run_client.step);
}
