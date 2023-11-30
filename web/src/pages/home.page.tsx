import React from "react";
import { Button, Container, Typography } from "@mui/material";
import { LineChart } from "@mui/x-charts/LineChart";
import { useNavigate } from "react-router-dom";
import ContactlessIcon from "@mui/icons-material/Contactless";
import PaidIcon from "@mui/icons-material/Paid";
import FeedIcon from "@mui/icons-material/Feed";
import SettingsIcon from "@mui/icons-material/Settings";
import EngineeringIcon from "@mui/icons-material/Engineering";

const HomePage = () => {
  const navigation = useNavigate();
  return (
    <Container maxWidth={false}>
      <div className="flex">
        {[
          {
            text: "OA Configs",
            link: "/oa-configs",
            icon: <SettingsIcon sx={{ fontSize: 80 }} />,
          },
          {
            text: "Integrator Configs",
            link: "/integrator-configs",
            icon: <EngineeringIcon sx={{ fontSize: 80 }} />,
          },
          {
            text: "Logs",
            link: "/logs",
            icon: <FeedIcon sx={{ fontSize: 80 }} />,
          },
          {
            text: "Online Authorisation Transactions",
            link: "/oa-transactions",
            icon: <ContactlessIcon sx={{ fontSize: 80 }} />,
          },
          {
            text: "Integrator Transactions",
            link: "/integrator-transactions",
            icon: <PaidIcon sx={{ fontSize: 80 }} />,
          },
        ].map(({ text, link, icon }, index) => (
          <Button className="flex-1 flex-col" onClick={() => navigation(link)}>
            {icon && icon}
            {text}
          </Button>
        ))}
      </div>

      <div className="h-8" />
      <Typography>[WIP]OA Transaction Error rate</Typography>
      <LineChart
        xAxis={[{ data: [1, 2, 3, 5, 8, 10] }]}
        series={[
          {
            data: [2, 5.5, 2, 8.5, 1.5, 5],
          },
        ]}
        height={300}
      />

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
