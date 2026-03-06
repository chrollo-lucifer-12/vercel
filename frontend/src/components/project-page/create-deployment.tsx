import { useCreateDeployment } from "@/hooks/use-deployment";
import { useEffect, useRef, useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";
import { Button } from "../ui/button";
import LogsDisplay from "./logs-display";
import { clientEnv } from "@/lib/env/client";
import { LogEvent } from "@/lib/types";

const CreateDeployment = ({ slug }: { slug: string }) => {
  const [logs, setLogs] = useState<LogEvent[]>([]);
  const [isDeploying, setIsDeploying] = useState(false);
  const [open, setOpen] = useState(false);

  const { mutate, data } = useCreateDeployment(slug);
  const eventSourceRef = useRef<EventSource | null>(null);

  useEffect(() => {
    if (!data?.deployment_id) return;

    if (eventSourceRef.current) {
      eventSourceRef.current.close();
      eventSourceRef.current = null;
    }

    const deploymentId = data.deployment_id;
    setLogs([]);
    setIsDeploying(true);

    const es = new EventSource(
      `http://localhost:9000/api/v1/deployment/logs/${deploymentId}`,
    );

    let cancelled = false;

    es.onopen = () => {
      console.log("SSE opened");
    };

    es.onmessage = (event) => {
      if (cancelled) return;
      if (event.data === "[DONE]") {
        es.close();
        setIsDeploying(false);
        return;
      }
      const raw = JSON.parse(event.data);

      const logEvent: LogEvent = {
        log: raw.log.message,
        created_at: raw.time,
        metadata: raw.log,
      };
      setLogs((prev) => [...prev, logEvent]);
    };

    es.onerror = () => {
      if (cancelled) return;
      if (es.readyState === EventSource.CLOSED) {
        setIsDeploying(false);
      }
    };

    eventSourceRef.current = es;

    return () => {
      cancelled = true;
      es.close();
      eventSourceRef.current = null;
    };
  }, [data?.deployment_id]);

  useEffect(() => {
    const handler = (e: BeforeUnloadEvent) => {
      if (!isDeploying) return;
      e.preventDefault();
      e.returnValue = "";
    };

    window.addEventListener("beforeunload", handler);
    return () => window.removeEventListener("beforeunload", handler);
  }, [isDeploying]);

  return (
    <Dialog
      open={open}
      onOpenChange={(val) => {
        if (isDeploying) return;
        setOpen(val);
      }}
    >
      <DialogTrigger asChild>
        <Button>Create Deployment</Button>
      </DialogTrigger>

      <DialogContent
        onEscapeKeyDown={(e) => {
          if (isDeploying) e.preventDefault();
        }}
        onPointerDownOutside={(e) => {
          if (isDeploying) e.preventDefault();
        }}
      >
        <DialogHeader>
          <DialogTitle>Create New Deployment</DialogTitle>
          <DialogDescription>
            Create a new deployment and see live logs.
          </DialogDescription>
        </DialogHeader>

        {logs.length > 0 && <LogsDisplay logs={logs} />}

        <Button
          onClick={() => {
            if (!slug) return;
            mutate();
          }}
        >
          Create
        </Button>
      </DialogContent>
    </Dialog>
  );
};

export default CreateDeployment;
