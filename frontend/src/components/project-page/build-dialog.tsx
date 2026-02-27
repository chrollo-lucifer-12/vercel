import { useGetDeployment } from "@/hooks/use-deployment";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";
import LogsDisplay from "./logs-display";
import { DotsThreeCircleIcon } from "@phosphor-icons/react";
import { Button } from "../ui/button";
import { ReactNode, useState } from "react";

const BuildDialog = ({
  deploymentId,
  children,
}: {
  deploymentId: string;
  children: ReactNode;
}) => {
  const [open, setOpen] = useState(false);
  const { data, isLoading } = useGetDeployment(deploymentId, open);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>
            Logs for deployment {data?.Deployment.created_at}
          </DialogTitle>
          <DialogDescription>
            Created At : {data?.Deployment.created_at}
            Status : {data?.Deployment.status}
          </DialogDescription>
        </DialogHeader>
        {isLoading ? (
          <p className="text-sm text-muted-foreground">Loading logs...</p>
        ) : (
          <LogsDisplay logs={data?.Logs!} />
        )}
      </DialogContent>
    </Dialog>
  );
};

export default BuildDialog;
