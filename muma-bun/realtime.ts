type Result<T, E> =
  | {
      ok: true;
      value: T;
    }
  | {
      ok: false;
      error: E;
    };

type Client = {
  clientId: string;
  controller: ReadableStreamDirectController;
};

class Session {
  public clients: Client[] = [];
  private data: string = "";

  constructor() {}

  public addClient = (controller: ReadableStreamDirectController): string => {
    const clientId = `client-${this.clients.length + 1}`;
    this.clients.push({
      clientId,
      controller,
    });

    return clientId;
  };

  public removeClient = (clientId: string) => {
    this.clients = this.clients.filter((c) => c.clientId !== clientId);
  };
}

export class Realtime {
  private sessions: Map<string, Session> = new Map();

  constructor() {}

  private createSession = (sessionId: string): Session => {
    const session = new Session();
    this.sessions.set(sessionId, session);
    console.log(session);
    return session;
  };

  private getSession = (sessionId: string): Result<Session, string> => {
    const session = this.sessions.get(sessionId);
    if (session == null) {
      return {
        ok: false,
        error: "No session found",
      };
    }

    return {
      ok: true,
      value: session,
    };
  };

  public emit = (sessionId: string, msg: string) => {
    const session = this.sessions.get(sessionId);

    console.log("asdasd", session);

    if (!session) {
      return;
    }

    for (const c of session.clients) {
      console.log("hereeee");
      c.controller.write(msg);
    }
  };

  private removeSession = (sessionId: string) => {};

  /**
   *
   * Creates a ndjson+jsonpatch streaming response
   *
   * @param req fetch Request instance
   * @param initialData initial response data
   * @param sessionId the session id to add to the list
   * @returns a streaming Response
   *
   * @example
   *
   * ```ts
   * // TODO: fill this example out
   * ```
   */
  public stream = <T>(req: Request, initialData: T, sessionId: string) => {
    let session: Session;
    const s = this.getSession(sessionId);

    if (s.ok) {
      session = s.value;
    } else {
      session = this.createSession(sessionId);
    }

    return new ReadableStream({
      type: "direct",
      pull(controller) {
        const clientId = session.addClient(controller);
        controller.write("asdf");
        if (req.signal.aborted) {
          console.log("client closed");
          controller.close();
          return;
        }
      },
    });
  };
}
