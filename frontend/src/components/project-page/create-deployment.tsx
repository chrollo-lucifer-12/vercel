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

const CreateDeployment = ({ slug }: { slug: string }) => {
  const [logs, setLogs] = useState<string[]>([]);
  const [isDeploying, setIsDeploying] = useState(false);
  const [open, setOpen] = useState(false);

  const { mutate, data } = useCreateDeployment(slug);
  const eventSourceRef = useRef<EventSource | null>(null);

  useEffect(() => {
    if (!data?.deployment_id) return;

    const deploymentId = data.deployment_id;

    const es = new EventSource(
      `http://localhost:9000/api/v1/deployment/logs/${deploymentId}`,
    );

    es.onmessage = (event) => {
      setIsDeploying(true);
      setLogs((prev) => [...prev, event.data]);
    };

    es.onerror = (err) => {
      console.error("SSE error:", err);
      es.close();
      setIsDeploying(false);
    };

    eventSourceRef.current = es;

    return () => es.close();
  }, [data]);

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
