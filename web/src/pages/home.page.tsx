import React from "react";
import {
  Card,
  CardContent,
  Container,
  Divider,
  Typography,
} from "@mui/material";
import { useQuery } from "react-query";
import { misc } from "../api/misc";

const HomePage = () => {
  const { data, refetch } = useQuery("misc", misc, {
    refetchInterval: 5000,
  });
  return (
    <Container>
      <div className="flex">
        <Card className="flex-1" sx={{ borderRadius: 10 }}>
          <CardContent>
            <Typography className="p-4" variant="h4">
              Total Entry
            </Typography>
            <Typography className="px-4" variant="h2">
              {data?.totalIn}
            </Typography>
          </CardContent>
        </Card>
        <div className="w-8" />
        <Card className="flex-1" sx={{ borderRadius: 10 }}>
          <CardContent>
            <Typography className="p-4" variant="h4">
              Total Exit
            </Typography>
            <Typography className="px-4" variant="h2">
              {data?.totalOut}
            </Typography>
          </CardContent>
        </Card>
      </div>
      <div className="h-8" />
      <Typography variant="h4">Integrator Status</Typography>
      <div className="h-2" />
      <div className="flex">
        {data?.integrators.map((integrator) => (
          <Card className="mr-8 flex">
            <Typography className="p-4" variant="h5">
              {integrator.integrator}
            </Typography>
            <Divider />
            <Typography
              className="p-4"
              style={{
                color: "white",
                backgroundColor:
                  integrator.status == "up" ? "#00afaa" : "#e4002b",
              }}
              variant="h5"
            >
              {integrator.status == "up" ? "Available" : "Error"}
            </Typography>
          </Card>
        ))}
      </div>
      <div className="h-8" />
      <Typography variant="h4">Snb Status</Typography>
      <div className="h-2" />
      <div className="flex">
        {data?.snb.map((snb) => (
          <Card className="mr-8 flex">
            <Typography className="p-4" variant="h5">
              {snb.facility}
            </Typography>
            <Divider />
            <Typography
              className="p-4"
              style={{
                color: "white",
                backgroundColor: snb.status == "up" ? "#00afaa" : "#e4002b",
              }}
              variant="h5"
            >
              {snb.status == "up" ? "Available" : "Error"}
            </Typography>
          </Card>
        ))}
      </div>
      <div className="h-8" />
      <Typography variant="h4">Last 10 Transaction</Typography>
      <div className="h-2" />
    </Container>
  );
};

export default HomePage;
