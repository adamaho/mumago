import { Realtime } from "./realtime";

function main() {
  const realtime = new Realtime();

  Bun.serve({
    port: 3000,
    fetch(req) {
      const url = new URL(req.url);

      switch (url.pathname) {
        case "/": {
          return new Response("hello world");
        }
        case "/bar": {
          const stream = realtime.stream(req, { data: [] }, "bar");
          return new Response(stream);
        }
        case "/foo": {
          realtime.emit("bar", "emit\n");
          return new Response("sent.");
        }
        default: {
          return new Response("default");
        }
      }

      // if (req.headers.get("X-Stream")) {
      //   const stream = new ReadableStream({
      //     type: "direct",
      //     async pull(controller) {
      //       while (!req.signal.aborted) {
      //         await delay();
      //         controller.write("world\n");
      //       }
      //     },
      //   });

      //   return new Response(stream);
      // }

      // return new Response("lol no stream");
    },
    tls: {
      cert: Bun.file("cert.pem"),
      key: Bun.file("key.pem"),
    },
  });

  console.log("Server started at https://localhost:3000");
}

main();
