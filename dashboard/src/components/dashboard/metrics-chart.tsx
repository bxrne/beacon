import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ResponsiveContainer, LineChart, Line, XAxis, YAxis, Tooltip, CartesianGrid } from "recharts";
import { memo } from "react";

interface MetricsChartProps {
  title: string;
  data: Array<{ timestamp: string; value: number }>;
  dataKey: string;
  color?: string;
}

// Memoize the tooltip component for better performance
const CustomTooltip = memo(({ active, payload }: any) => {
  if (!active || !payload?.length) return null;

  return (
    <div className="rounded-lg border bg-background p-2 shadow-sm">
      <div className="grid grid-cols-2 gap-2">
        <div className="flex flex-col">
          <span className="text-[0.70rem] uppercase text-muted-foreground">
            Time
          </span>
          <span className="font-bold text-muted-foreground">
            {new Date(payload[0].payload.timestamp).toLocaleTimeString()}
          </span>
        </div>
        <div className="flex flex-col">
          <span className="text-[0.70rem] uppercase text-muted-foreground">
            Value
          </span>
          <span className="font-bold">
            {payload[0].value}
          </span>
        </div>
      </div>
    </div>
  );
});

CustomTooltip.displayName = "CustomTooltip";

// Memoize the entire chart component
export const MetricsChart = memo(({ 
  title, 
  data, 
  dataKey, 
  color = "hsl(var(--chart-1))" 
}: MetricsChartProps) => {
  const formatXAxis = (value: string) => new Date(value).toLocaleTimeString();

  return (
    <Card className="col-span-4">
      <CardHeader>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent className="h-[300px]">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart data={data} margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
            <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
            <XAxis
              dataKey="timestamp"
              tickFormatter={formatXAxis}
              className="text-sm text-muted-foreground"
              padding={{ left: 10, right: 10 }}
              tick={{ fontSize: 12 }}
            />
            <YAxis 
              className="text-sm text-muted-foreground"
              tick={{ fontSize: 12 }}
              width={40}
            />
            <Tooltip content={CustomTooltip} />
            <Line
              type="monotone"
              dataKey={dataKey}
              stroke={color}
              strokeWidth={2}
              dot={false}
              activeDot={{ r: 4, strokeWidth: 0 }}
              isAnimationActive={false}
            />
          </LineChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
});

MetricsChart.displayName = "MetricsChart";