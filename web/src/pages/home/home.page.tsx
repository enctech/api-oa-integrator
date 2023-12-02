import React from "react";
import { Container, Typography } from "@mui/material";
import { LineChart } from "@mui/x-charts/LineChart";
import OAErrorRates from "./oa-error-rate";

const HomePage = () => {
  return (
    <Container maxWidth={false}>
      <div className="h-8" />
      <OAErrorRates />

      <div className="h-8" />
      <Typography>[WIP]Integrator Error rate</Typography>
      <LineChart
        xAxis={[{ data: [1, 2, 3, 5, 8, 10] }]}
        series={[
          {
            data: [2, 5.5, 2, 8.5, 1.5, 5],
          },
        ]}
        height={300}
      />
    </Container>
  );
};

export default HomePage;
