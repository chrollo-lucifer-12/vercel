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
import { CalendarBlankIcon, DotsThreeCircleIcon } from "@phosphor-icons/react";
import { ReactNode, useState } from "react";
import { Skeleton } from "../ui/skeleton";
import { Badge } from "../ui/badge";
import { formatDate } from "@/lib/utils";

const BuildDialogSkeleton = () => {
  const rows = Array.from({ length: 8 });

  return (
    <div className="flex flex-col gap-2">
      <Skeleton className="h-6 w-1/3 rounded-md" />

      <Skeleton className="h-4 w-1/4 rounded-md" />
      <Skeleton className="h-4 w-1/4 rounded-md" />

      <div className="mt-2 flex flex-col gap-1">
        {rows.map((_, idx) => (
          <Skeleton key={idx} className="h-3 w-full rounded-md" />
        ))}
      </div>
    </div>
  );
};

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
          <DialogDescription className="flex flex-row gap-2 items-center">
            <Badge variant={"secondary"}>
              <CalendarBlankIcon data-icon="inline-start" />
              {formatDate(data?.Deployment.created_at!)}
            </Badge>
            <Badge
              variant={`${data?.Deployment.status === "FAILED" ? "destructive" : "secondary"}`}
            >
              {data?.Deployment.status}
            </Badge>
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
