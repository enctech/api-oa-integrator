import React, { useEffect, useRef } from "react";
import { Controller, SubmitHandler, useForm } from "react-hook-form";
import { useDebounce } from "@uidotdev/usehooks";
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Container,
  Grid,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TablePagination,
  TableRow,
  TextField,
  Typography,
} from "@mui/material";
import { useSearchParams } from "react-router-dom";
import { useQuery } from "react-query";
import { getOATransactions } from "../api/transactions";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import { DateTimePicker } from "@mui/x-date-pickers";
import moment from "moment";
import dayjs, { Dayjs } from "dayjs";

const statusMapper: Map<string, string> = new Map([
  ["identification_entry_start", "User entry verification"],
  ["identification_entry_done", "User entry verification done"],
  ["identification_entry_error", "User entry verification failed"],
  ["leave_loop_entry_done", "User entered"],
  ["payment_exit_error", "Payment failed"],
]);

interface FormData {
  startAt?: Dayjs;
  endAt?: Dayjs;
  exitLane: string;
  entryLane: string;
  lpn: string;
  facility: string;
  jobId: string;
}

const OATransactionsPage: React.FC = () => {
  const perPagesDefault = useRef([100, 500, 1000]);
  const [currentQueryParameters, setSearchParams] = useSearchParams();
  const newParams = useRef(new URLSearchParams());
  const { control, register, handleSubmit, watch } = useForm<FormData>({
    defaultValues: {
      lpn: currentQueryParameters.get("lpn") || "",
      jobId: currentQueryParameters.get("jobId") || "",
      entryLane: currentQueryParameters.get("entryLane") || "",
      facility: currentQueryParameters.get("facility") || "",
      exitLane: currentQueryParameters.get("exitLane") || "",
      startAt: dayjs(
        moment(currentQueryParameters.get("startAt")).local().toDate(),
      ),
      endAt: dayjs(
        moment(currentQueryParameters.get("endAt")).local().toDate(),
      ),
    },
  });

  const debouncedSearchTerm = useDebounce(currentQueryParameters, 300);

  const { data, refetch } = useQuery(
    "getOATransactions",
    () =>
      getOATransactions({
        page: +(currentQueryParameters.get("page") || "0"),
        perPage: +(
          currentQueryParameters.get("perPage") ||
          `${perPagesDefault.current[0]}`
        ),
        startAt: currentQueryParameters.get("startAt")
          ? moment(currentQueryParameters.get("startAt")).utc().toDate()
          : undefined,
        endAt: currentQueryParameters.get("endAt")
          ? moment(currentQueryParameters.get("endAt")).utc().toDate()
          : undefined,
        exitLane: currentQueryParameters.get("exitLane") || undefined,
        entryLane: currentQueryParameters.get("entryLane") || undefined,
        lpn: currentQueryParameters.get("lpn") || undefined,
        facility: currentQueryParameters.get("facility") || undefined,
        jobId: currentQueryParameters.get("jobId") || undefined,
      }),
    {
      refetchInterval: 5000,
    },
  );

  useEffect(() => {
    if (debouncedSearchTerm) {
      refetch();
    }
  }, [debouncedSearchTerm]);

  useEffect(() => {
    newParams.current.set("page", "0");
    newParams.current.set("perPage", "100");
    if (currentQueryParameters.get("startAt")) {
      newParams.current.set(
        "startAt",
        moment(currentQueryParameters.get("startAt")).utc().toISOString(),
      );
    }
    if (currentQueryParameters.get("endAt")) {
      newParams.current.set(
        "endAt",
        moment(currentQueryParameters.get("endAt")).utc().toISOString(),
      );
    }
    newParams.current.set(
      "exitLane",
      currentQueryParameters.get("exitLane") || "",
    );
    newParams.current.set(
      "entryLane",
      currentQueryParameters.get("entryLane") || "",
    );
    newParams.current.set("lpn", currentQueryParameters.get("lpn") || "");
    newParams.current.set(
      "facility",
      currentQueryParameters.get("facility") || "",
    );
    newParams.current.set("jobId", currentQueryParameters.get("jobId") || "");

    setSearchParams(newParams.current);
  }, []);

  const onSubmit: SubmitHandler<FormData> = (formData) => {
    if (formData.startAt && formData.startAt.isValid()) {
      newParams.current.set("startAt", formData.startAt.toISOString());
    }
    if (formData.endAt && formData.endAt.isValid()) {
      newParams.current.set("endAt", formData.endAt.toISOString());
    }
    newParams.current.set("exitLane", formData.exitLane || "");
    newParams.current.set("entryLane", formData.entryLane || "");
    newParams.current.set("lpn", formData.lpn || "");
    newParams.current.set("facility", formData.facility || "");
    newParams.current.set("jobId", formData.jobId || "");

    setSearchParams(newParams.current);
  };

  // Watch for changes in form fields
  watch([
    "startAt",
    "endAt",
    "exitLane",
    "entryLane",
    "lpn",
    "facility",
    "jobId",
  ]);

  // Function to call API when there are changes in any field
  const handleFieldChange = () => {
    handleSubmit(onSubmit)();
  };

  return (
    <Container maxWidth={false}>
      <Accordion>
        <AccordionSummary
          expandIcon={<ExpandMoreIcon />}
          aria-controls="panel1a-content"
          id="panel1a-header"
        >
          <Typography>Filter</Typography>
        </AccordionSummary>
        <AccordionDetails>
          <form>
            <Grid container spacing={2}>
              <Grid item xs={3}>
                <Controller
                  name="startAt"
                  control={control}
                  render={({ field }) => (
                    <DateTimePicker
                      className="flex-1 w-full"
                      {...field}
                      onChange={(_) => {}}
                      label="Start Date"
                      onAccept={(value) => {
                        if (!value) return;
                        field.onChange(value);
                        handleFieldChange();
                      }}
                    />
                  )}
                />
              </Grid>
              <Grid item xs={3}>
                <Controller
                  name="endAt"
                  control={control}
                  render={({ field }) => (
                    <DateTimePicker
                      className="flex-1 w-full"
                      {...field}
                      format={"DD/MM/YYYY hh:mm"}
                      label="End Date"
                      onAccept={(value) => {
                        if (!value) return;
                        field.onChange(value);
                        handleFieldChange();
                      }}
                    />
                  )}
                />
              </Grid>

              {[
                {
                  type: "exitLane",
                  label: "Exit Lane",
                },
                {
                  type: "entryLane",
                  label: "Entry Lane",
                },
                {
                  type: "lpn",
                  label: "License Plate Number",
                },
                {
                  type: "facility",
                  label: "Facility",
                },
                {
                  type: "jobId",
                  label: "Job",
                },
              ].map((fieldName) => (
                <Grid item xs={6} key={fieldName.type}>
                  <TextField
                    variant="outlined"
                    fullWidth
                    label={fieldName.label}
                    {...register(fieldName.type as any, {
                      onChange: (_) => handleFieldChange(),
                    })}
                  />
                </Grid>
              ))}
            </Grid>
          </form>
        </AccordionDetails>
      </Accordion>

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
                  {moment(row.createdAt).local().format("DD/MM/yyyy hh:mm:ss")}
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
        <TablePagination
          rowsPerPageOptions={perPagesDefault.current}
          component="div"
          count={data?.data?.length || 0}
          rowsPerPage={parseInt(
            currentQueryParameters.get("perPage") || "0",
            10,
          )}
          page={parseInt(currentQueryParameters.get("page") || "0", 10)}
          onPageChange={(_: any, newPage: number) => {
            newParams.current.set("page", newPage.toString());
            setSearchParams(newParams.current);
          }}
          onRowsPerPageChange={(event: { target: { value: string } }) => {
            newParams.current.set("perPage", event.target.value);
            setSearchParams(newParams.current);
          }}
        />
      </TableContainer>
    </Container>
  );
};

export default OATransactionsPage;
