import { useCreateDeployment } from "@/hooks/use-deployment";
import { Button } from "../ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";
import { useEffect, useRef, useState } from "react";
import LogsDisplay from "./logs-display";

const CreateDeployment = ({ slug }: { slug: string }) => {
  console.log(slug);
  const [logs, setLogs] = useState<string[]>([]);
  const { mutate, data } = useCreateDeployment(slug);
  const eventSourceRef = useRef<EventSource | null>(null);

  useEffect(() => {
    if (!data?.deployment_id) return;

    const deploymentId = data.deployment_id;

    const es = new EventSource(
      `http://localhost:9000/api/v1/deployment/logs/${deploymentId}`,
    );

    es.onmessage = (event) => {
      console.log(event.data);
      setLogs((prev) => [...prev, event.data]);
    };

    es.onerror = (err) => {
      console.error("SSE error:", err);
      es.close();
    };

    eventSourceRef.current = es;

    return () => {
      es.close();
    };
  }, [data]);

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button>Create Deployment</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Crearte New Deployment</DialogTitle>
          <DialogDescription>
            Create a new deployment for your project and see live logs.
          </DialogDescription>
        </DialogHeader>
        {/*<LogsDisplay logs={logs} />*/}
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
