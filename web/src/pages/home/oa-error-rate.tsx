import { Typography } from "@mui/material";
import { LineChart } from "@mui/x-charts/LineChart";
import React from "react";
import moment from "moment";

export default function OAErrorRates() {
  return (
    <div>
      <Typography>[WIP]OA Transaction Error rate</Typography>
      <LineChart
        xAxis={[
          {
            data: [1, 2, 3, 5, 8, 10],
            min: moment().startOf("day").toDate(),
            max: moment().endOf("day").toDate(),
          },
        ]}
        series={[
          {
            data: [2, 5.5, 2, 8.5, 1.5, 5],
          },
        ]}
        height={300}
      />
    </div>
  );
}
