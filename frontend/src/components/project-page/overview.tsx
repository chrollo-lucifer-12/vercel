import { CalendarBlankIcon, LinkIcon } from "@phosphor-icons/react";
import { Badge } from "../ui/badge";
import { Card, CardContent } from "../ui/card";
import { TabsContent } from "../ui/tabs";
import { LogEvent } from "@/lib/types";
import LogsDisplay from "./logs-display";
import { formatDate } from "@/lib/utils";

const Overview = ({
  logs,
  subDomain,
  createdAt,
}: {
  logs: LogEvent[];
  subDomain: string;
  createdAt: string;
}) => {
  return (
    <TabsContent value="overview" className="mt-6 flex flex-col gap-2">
      <p className="text-2xl font-semibold">Current Deployment</p>
      <Card>
        <CardContent className="flex flex-col gap-6">
          <div className="flex justify-between h-48 items-center">
            <iframe
              src="https://vercel.com"
              className="w-[90%] h-full border rounded"
            />
            <div className="w-[60%] pl-4 flex flex-col justify-center gap-2">
              <Badge variant={"secondary"}>
                <CalendarBlankIcon data-icon="inline-start" />
                {formatDate(createdAt)}
              </Badge>
              <Badge variant={"secondary"}>
                <LinkIcon data-icon="inline-start" />
                {subDomain}
              </Badge>
            </div>
          </div>
          <div>
            <h1 className="text-lg">Build Logs</h1>
            <LogsDisplay logs={logs} />
          </div>
        </CardContent>
      </Card>
    </TabsContent>
  );
};

export default Overview;
