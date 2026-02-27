import { useAnalytics } from "@/hooks/use-project";
import { TabsContent } from "../ui/tabs";
import { useState, useMemo } from "react";
import { format } from "date-fns";

type DateRange = { from: Date | undefined; to?: Date | undefined };
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import {
  Area,
  AreaChart,
  Bar,
  BarChart,
  CartesianGrid,
  XAxis,
  YAxis,
} from "recharts";
import { WebsiteAnalytics } from "@/lib/types";
import { CalendarIcon, XIcon } from "@phosphor-icons/react";

const responseTimeConfig = {
  response_time_ms: {
    label: "Avg Response Time (ms)",
    color: "hsl(var(--chart-1))",
  },
} satisfies ChartConfig;

const requestCountConfig = {
  count: {
    label: "Requests",
    color: "hsl(var(--chart-2))",
  },
} satisfies ChartConfig;

function formatDate(dateStr: string) {
  const d = new Date(dateStr);
  return d.toLocaleDateString("en-US", { month: "short", day: "numeric" });
}

function groupByDate(data: WebsiteAnalytics[]) {
  const map = new Map<
    string,
    { count: number; totalResponseTime: number; responseCount: number }
  >();

  for (const item of data) {
    const key = formatDate(item.created_at);
    const existing = map.get(key) ?? {
      count: 0,
      totalResponseTime: 0,
      responseCount: 0,
    };
    existing.count += 1;
    if (item.response_time_ms != null) {
      existing.totalResponseTime += item.response_time_ms;
      existing.responseCount += 1;
    }
    map.set(key, existing);
  }

  return Array.from(map.entries())
    .map(([date, { count, totalResponseTime, responseCount }]) => ({
      date,
      count,
      response_time_ms:
        responseCount > 0
          ? Math.round(totalResponseTime / responseCount)
          : null,
    }))
    .sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());
}

const PRESETS = [
  { label: "Last 7 days", days: 7 },
  { label: "Last 30 days", days: 30 },
  { label: "Last 90 days", days: 90 },
];

const Analytics = ({ subDomain }: { subDomain: string }) => {
  const [dateRange, setDateRange] = useState<DateRange | undefined>(undefined);
  const from = dateRange?.from ?? null;
  const to = dateRange?.to ?? null;
  const { data, isLoading } = useAnalytics(subDomain, from, to);

  const chartData = useMemo(() => {
    if (!data) return [];
    return groupByDate(data);
  }, [data]);

  function applyPreset(days: number) {
    const to = new Date();
    const from = new Date();
    from.setDate(from.getDate() - days + 1);
    setDateRange({ from, to });
  }

  function clearRange() {
    setDateRange(undefined);
  }

  const rangeLabel =
    dateRange?.from && dateRange?.to
      ? `${format(dateRange.from, "MMM d, yyyy")} – ${format(dateRange.to, "MMM d, yyyy")}`
      : dateRange?.from
        ? `${format(dateRange.from, "MMM d, yyyy")} – ...`
        : "Pick a date range";

  return (
    <TabsContent value="analytics" className="mt-6 flex flex-col gap-6">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <p className="text-2xl font-semibold">Analytics</p>

        <div className="flex flex-wrap items-center gap-2">
          {PRESETS.map((preset) => (
            <Button
              key={preset.label}
              variant="outline"
              size="sm"
              onClick={() => applyPreset(preset.days)}
            >
              {preset.label}
            </Button>
          ))}

          <Popover>
            <PopoverTrigger asChild>
              <Button variant="outline" size="sm" className="gap-2">
                <CalendarIcon className="h-4 w-4" />
                <span className="hidden sm:inline">{rangeLabel}</span>
                <span className="sm:hidden">Custom</span>
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-auto p-0" align="end">
              <Calendar
                mode="range"
                selected={dateRange}
                onSelect={setDateRange}
                numberOfMonths={2}
              />
            </PopoverContent>
          </Popover>

          {dateRange && (
            <Button
              variant="ghost"
              size="sm"
              onClick={clearRange}
              className="gap-1 px-2"
            >
              <XIcon className="h-4 w-4" />
              Clear
            </Button>
          )}
        </div>
      </div>

      {isLoading && (
        <p className="text-muted-foreground text-sm">Loading analytics...</p>
      )}

      {!isLoading && chartData.length === 0 && (
        <p className="text-muted-foreground text-sm">
          No analytics data available.
        </p>
      )}

      {!isLoading && chartData.length > 0 && (
        <>
          {/* Request Count Chart */}
          <Card>
            <CardHeader>
              <CardTitle>Request Volume</CardTitle>
              <CardDescription>Number of requests over time</CardDescription>
            </CardHeader>
            <CardContent>
              <ChartContainer
                config={requestCountConfig}
                className="h-[250px] w-full"
              >
                <BarChart data={chartData} margin={{ left: 0, right: 8 }}>
                  <CartesianGrid vertical={false} />
                  <XAxis
                    dataKey="date"
                    tickLine={false}
                    axisLine={false}
                    tickMargin={8}
                    tick={{ fontSize: 12 }}
                  />
                  <YAxis
                    tickLine={false}
                    axisLine={false}
                    tickMargin={8}
                    tick={{ fontSize: 12 }}
                    allowDecimals={false}
                  />
                  <ChartTooltip content={<ChartTooltipContent />} />
                  <Bar
                    dataKey="count"
                    fill="var(--color-count)"
                    radius={[4, 4, 0, 0]}
                  />
                </BarChart>
              </ChartContainer>
            </CardContent>
          </Card>

          {/* Response Time Chart */}
          <Card>
            <CardHeader>
              <CardTitle>Average Response Time</CardTitle>
              <CardDescription>
                Average response time in milliseconds over time
              </CardDescription>
            </CardHeader>
            <CardContent>
              <ChartContainer
                config={responseTimeConfig}
                className="h-[250px] w-full"
              >
                <AreaChart data={chartData} margin={{ left: 0, right: 8 }}>
                  <defs>
                    <linearGradient
                      id="fillResponseTime"
                      x1="0"
                      y1="0"
                      x2="0"
                      y2="1"
                    >
                      <stop
                        offset="5%"
                        stopColor="var(--color-response_time_ms)"
                        stopOpacity={0.3}
                      />
                      <stop
                        offset="95%"
                        stopColor="var(--color-response_time_ms)"
                        stopOpacity={0}
                      />
                    </linearGradient>
                  </defs>
                  <CartesianGrid vertical={false} />
                  <XAxis
                    dataKey="date"
                    tickLine={false}
                    axisLine={false}
                    tickMargin={8}
                    tick={{ fontSize: 12 }}
                  />
                  <YAxis
                    tickLine={false}
                    axisLine={false}
                    tickMargin={8}
                    tick={{ fontSize: 12 }}
                    unit="ms"
                  />
                  <ChartTooltip content={<ChartTooltipContent />} />
                  <Area
                    dataKey="response_time_ms"
                    type="monotone"
                    stroke="var(--color-response_time_ms)"
                    fill="url(#fillResponseTime)"
                    strokeWidth={2}
                    dot={false}
                    connectNulls
                  />
                </AreaChart>
              </ChartContainer>
            </CardContent>
          </Card>
        </>
      )}
    </TabsContent>
  );
};

export default Analytics;
