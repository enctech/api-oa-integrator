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
import TimeToLeaveIcon from "@mui/icons-material/TimeToLeave";
import NoCrashIcon from "@mui/icons-material/NoCrash";
import PaidIcon from "@mui/icons-material/Paid";
import { useQuery } from "react-query";
import { integratorStatus, misc } from "../api/misc";
import { getLatestOATransactions } from "../api/transactions";
import moment from "moment";
import { statusMapper } from "./oa-transactions.page";

const HomePage = () => {
  const { data } = useQuery("misc", misc, {
    refetchInterval: 1000 * 5,
  });
  const { data: integratorData } = useQuery(
    "integratorStatus",
    integratorStatus,
    {
      refetchInterval: 1000 * 60,
    },
  );
  return (
    <Container>
      <div className="flex">
        <Card
          sx={{
            borderRadius: 4,
            flex: 1,
            backgroundColor: "#fff0bf80",
            display: "flex",
            alignItems: "center",
          }}
          elevation={10}
        >
          <CardContent>
            <TimeToLeaveIcon
              sx={{
                fontSize: 40,
              }}
            />
            <div className="h-2" />
            <Typography variant="h5" fontWeight="bold">
              {data?.totalIn}
            </Typography>
            <Typography variant="body1">Total Entry</Typography>
          </CardContent>
        </Card>
        <div className="w-8" />
        <Card
          sx={{
            borderRadius: 4,
            flex: 1,
            backgroundColor: "#ffe27f80",
            display: "flex",
            alignItems: "center",
          }}
          elevation={10}
        >
          <CardContent>
            <NoCrashIcon
              sx={{
                fontSize: 40,
              }}
            />
            <div className="h-2" />
            <Typography variant="h5">{data?.totalOut}</Typography>
            <Typography variant="body1">Total Exit</Typography>
          </CardContent>
        </Card>
        <div className="w-8" />
        <Card
          sx={{
            borderRadius: 4,
            flex: 1,
            backgroundColor: "#00afaa80",
            display: "flex",
            alignItems: "center",
          }}
          elevation={10}
        >
          <CardContent>
            <PaidIcon
              sx={{
                fontSize: 40,
              }}
            />
            <div className="h-2" />
            <Typography variant="h5">
              RM {(data?.totalPayment || 0.0).toFixed(2)}
            </Typography>
            <Typography variant="body1">Total Payment</Typography>
          </CardContent>
        </Card>
      </div>
      <div className="h-8" />
      <div className="flex">
        <Status
          title="3rd parties Status"
          data={
            integratorData?.integrators?.map((x) => ({
              info: x.integrator,
              status: x.status,
            })) || []
          }
          partialAvailableMessage="Partial 3rd parties available"
          fullyAvailableMessage="All 3rd parties available"
          noAvailableMessage="All 3rd parties unavailable"
        />
        <div className="w-4" />
        <Status
          title="Snb Status"
          data={
            integratorData?.snb?.map((x) => ({
              info: x.facility,
              status: x.status,
            })) || []
          }
          partialAvailableMessage="Partial SnB system available"
          fullyAvailableMessage="All SnB system available"
          noAvailableMessage="All SnB system unavailable"
        />
      </div>
      <div className="h-8" />
      <Typography variant="h5">Last 10 Transaction</Typography>
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
        page: 0,
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
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

const Status = (input: {
  data: { info: string; status: string }[];
  title: string;
  partialAvailableMessage: string;
  fullyAvailableMessage: string;
  noAvailableMessage: string;
}) => {
  const message = () => {
    if (input.data.every((i) => i.status === "up")) {
      return input.fullyAvailableMessage;
    }

    if (input.data.every((i) => i.status !== "up")) {
      return input.noAvailableMessage;
    }

    return input.partialAvailableMessage;
  };
  return (
    <Card
      elevation={10}
      sx={{
        flex: 1,
        borderRadius: 4,
      }}
    >
      <CardContent>
        <Typography variant="h5">{input.title}</Typography>
        <div className="h-2" />
        <div>{message()}</div>
        <div>
          {input.data.map((snb) => (
            <>
              <div className="mt-2 flex items-center">
                <Typography className="p-2 flex-1" variant="body1">
                  {snb.info}
                </Typography>
                <Typography
                  className="p-2"
                  style={{
                    width: 100,
                    fontWeight: "bold",
                    textAlign: "right",
                    color: snb.status == "up" ? "#00afaa" : "#e4002b",
                  }}
                  variant="body1"
                >
                  {snb.status == "up" ? "UP" : "DOWN"}
                </Typography>
                <div className="w-2" />
                {snb.status == "up" ? (
                  <div className="w-5 h-5 rounded-full inline-flex items-center justify-center bg-[#00afaa]" />
                ) : (
                  <div className="w-5 h-5 rounded-full inline-flex items-center justify-center bg-[#e4002b]" />
                )}
              </div>
              <Divider />
            </>
          ))}
        </div>
      </CardContent>
    </Card>
  );
};

export default HomePage;
