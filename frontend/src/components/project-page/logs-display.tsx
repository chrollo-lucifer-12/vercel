import { LogEvent } from "@/lib/types";
import { formatDate } from "@/lib/utils";

const LogsDisplay = ({ logs }: { logs: LogEvent[] }) => {
  return (
    <div className="bg-black text-green-400 font-mono text-sm p-4 rounded-lg h-96 overflow-y-auto">
      {logs
        ?.filter((log) => {
          return log.log.length > 0;
        })
        .map((log, i) => (
          <div key={i} className="whitespace-pre-wrap">
            {formatDate(log.created_at)} {log.log}
          </div>
        ))}
    </div>
  );
};

export default LogsDisplay;
