import React from 'react';
import type { Metric } from '../../types/api';
import { format } from 'date-fns';

interface MetricCardProps {
  metric: Metric;
}

export const MetricCard: React.FC<MetricCardProps> = ({ metric }) => {
  const isQualitative = !Number(metric.value);

  return (
    <div className="bg-white rounded-xl p-4 shadow-sm hover:shadow-md transition-shadow">
      <div className="flex justify-between items-start mb-2">
        <h3 className="text-sm font-medium text-gray-500 uppercase">{metric.type}</h3>
        <span className="text-xs text-gray-400">
          {format(new Date(metric.recorded_at), 'HH:mm')}
        </span>
      </div>
      
      <div className="mt-2">
        {isQualitative ? (
          <span className="text-lg font-semibold text-gray-900">{metric.value}</span>
        ) : (
          <div className="flex items-baseline">
            <span className="text-2xl font-bold text-gray-900">{metric.value}</span>
            <span className="ml-1 text-sm text-gray-500">{metric.unit}</span>
          </div>
        )}
      </div>
    </div>
  );
};