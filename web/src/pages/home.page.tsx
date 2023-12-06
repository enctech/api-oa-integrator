import React from "react";
import {
  Card,
  CardContent,
  Container,
  Divider,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from "@mui/material";
import { useQuery } from "react-query";
import { misc } from "../api/misc";
import { getLatestOATransactions } from "../api/transactions";
import moment from "moment";
import { statusMapper } from "./oa-transactions.page";

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
      <LatestTransactions />
    </Container>
  );
};

const LatestTransactions = () => {
  const { data } = useQuery(
    ["getLatestOATransactions"],
    () =>
      getLatestOATransactions({
        startAt: moment().startOf("day").utc().toDate(),
        endAt: moment().endOf("day").utc().toDate(),
        page: 1,
        perPage: 10,
      }),
    {
      refetchInterval: 5000,
    },
  );
  return (
    <TableContainer component={Paper} className="mt-4">
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>Created At</TableCell>
            <TableCell>License Plate Number</TableCell>
            <TableCell>Entry Lane</TableCell>
            <TableCell>Exit Lane</TableCell>
            <TableCell>
              <div className="w-32">Status</div>
            </TableCell>
            <TableCell>
              <div className="w-96">Error</div>
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {data?.data?.map((row) => (
            <TableRow key={row.id}>
              <TableCell>
                {moment(row.createdAt).local().format("DD/MM/yyyy hh:mm:ss A")}
              </TableCell>
              <TableCell>{row.lpn}</TableCell>
              <TableCell>{row.entryLane}</TableCell>
              <TableCell>{row.exitLane || "-"}</TableCell>
              <TableCell>
                <div className="w-32">
                  {statusMapper.get(row.extra.steps) || row.extra.steps}
                </div>
              </TableCell>
              <TableCell
                style={{
                  width: "30px",
                  whiteSpace: "normal",
                  wordWrap: "break-word",
                }}
              >
                <div className="w-[40rem]">{row.extra.error}</div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

export default HomePage;
