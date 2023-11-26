import React, { useEffect, useRef, useState } from "react";
import {
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TablePagination,
  TableRow,
  TextField,
} from "@mui/material";
import { useQuery } from "react-query";
import { getOALogs } from "../api/transactions";
import { JsonViewer } from "@textea/json-viewer";
import { DateTimePicker } from "@mui/x-date-pickers";
import dayjs from "dayjs";
import { useSearchParams } from "react-router-dom";

const LogsPage = () => {
  const [currentQueryParameters, setSearchParams] = useSearchParams();

  const perPagesDefault = useRef([100, 500, 1000]);

  const [rowsPerPage, setRowsPerPage] = useState(perPagesDefault.current[0]);
  const [page, setPage] = useState(0);

  const [beforeTime, setBeforeTime] = useState(
    currentQueryParameters.get("before")
      ? new Date(currentQueryParameters.get("before") as string)
      : new Date(),
  );
  const [afterTime, setAfterTime] = useState(
    currentQueryParameters.get("after")
      ? new Date(currentQueryParameters.get("after") as string)
      : new Date(Date.now() - 1000 * 60 * 5),
  );
  const [message, setMessage] = useState("");
  const [field, setField] = useState("");

  const { data, refetch } = useQuery("getOALogs", () =>
    getOALogs({
      after: new Date(currentQueryParameters.get("after") || afterTime),
      before: new Date(currentQueryParameters.get("before") || beforeTime),
      message: currentQueryParameters.get("message"),
      field: currentQueryParameters.get("field"),
      page: +(currentQueryParameters.get("page") || "0"),
      perPage: +(
        currentQueryParameters.get("perPage") || `${perPagesDefault.current[0]}`
      ),
    }),
  );

  useEffect(() => {
    const newParams = new URLSearchParams();

    newParams.set("before", beforeTime.toISOString());
    newParams.set("after", afterTime.toISOString());
    newParams.set("message", message);
    newParams.set("field", field);
    newParams.set("page", `${page}`);
    newParams.set("perPage", `${rowsPerPage}`);

    setSearchParams(newParams);
  }, [beforeTime, afterTime, message, field, page, rowsPerPage]);

  useEffect(() => {
    refetch().then();
  }, [currentQueryParameters]);

  const handleChangePage = (_: any, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: { target: { value: string } }) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  return (
    <div className="p-4 px-8">
      <div className={"flex"}>
        <div className={"py-4"}>
          <DateTimePicker
            label="After"
            format={"DD/MM/YYYY hh:mm"}
            defaultValue={dayjs(
              new Date(currentQueryParameters.get("after") || afterTime),
            )}
            onAccept={(value) => {
              if (value) setAfterTime(value.toDate());
            }}
          />
        </div>
        <div className={"p-4"}>
          <DateTimePicker
            label="Before"
            format={"DD/MM/YYYY hh:mm"}
            defaultValue={dayjs(
              new Date(currentQueryParameters.get("before") || beforeTime),
            )}
            onAccept={(value) => {
              if (value) setBeforeTime(value.toDate());
            }}
          />
        </div>
        <div className={"p-4"}>
          <TextField
            label="Message"
            value={message}
            onChange={(e) => setMessage(e.target.value)}
          />
        </div>
        <div className={"p-4"}>
          <TextField
            label="Field filter"
            value={field}
            onChange={(e) => setField(e.target.value)}
          />
        </div>
      </div>
      <TableContainer component={Paper} className="mt-4">
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Level</TableCell>
              <TableCell>Message</TableCell>
              <TableCell>Fields</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data?.data?.map((row) => (
              <TableRow key={row.id}>
                <TableCell>{row.level}</TableCell>
                <TableCell>{row.message}</TableCell>
                <TableCell>
                  <JsonViewer
                    value={row.fields}
                    defaultInspectControl={(_, __) => false}
                    collapseStringsAfterLength={false}
                  />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        <TablePagination
          rowsPerPageOptions={perPagesDefault.current}
          component="div"
          count={data?.metadata?.totalData || 0}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={handleChangePage}
          onRowsPerPageChange={handleChangeRowsPerPage}
        />
      </TableContainer>
    </div>
  );
};

export default LogsPage;
