import React, { useEffect, useRef } from "react";
import { Controller, SubmitHandler, useForm } from "react-hook-form";
import { useDebounce } from "@uidotdev/usehooks";
import { JsonViewer } from "@textea/json-viewer";
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
import { getIntegratorTransactions } from "../api/transactions";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import { ClearIcon, DateTimePicker } from "@mui/x-date-pickers";
import moment from "moment";
import dayjs, { Dayjs } from "dayjs";
import IconButton from "@mui/material/IconButton";

interface FormData {
  startAt?: Dayjs | null;
  endAt?: Dayjs | null;
  status: string;
  lpn: string;
  integratorName: string;
}

const IntegratorTransactionsPage: React.FC = () => {
  const perPagesDefault = useRef([100, 500, 1000]);
  const [currentQueryParameters, setSearchParams] = useSearchParams();
  const newParams = useRef(new URLSearchParams());
  const { control, register, handleSubmit, watch, setValue } =
    useForm<FormData>({
      defaultValues: {
        lpn: currentQueryParameters.get("lpn") || "",
        status: currentQueryParameters.get("status") || "",
        integratorName: currentQueryParameters.get("integratorName") || "",
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
      getIntegratorTransactions({
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
        lpn: currentQueryParameters.get("lpn") || undefined,
        status: currentQueryParameters.get("lpn") || undefined,
        integratorName:
          currentQueryParameters.get("integratorName") || undefined,
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

    newParams.current.set("status", formData.status || "");
    newParams.current.set("lpn", formData.lpn || "");
    newParams.current.set("integratorName", formData.integratorName || "");

    setSearchParams(newParams.current);
  };

  // Watch for changes in form fields
  watch(["startAt", "endAt", "lpn", "status", "integratorName"]);

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
                    <div className="flex">
                      <DateTimePicker
                        format={"DD/MM/YYYY hh:mm"}
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
                      <IconButton
                        onClick={() => {
                          setValue("startAt", null);
                          newParams.current.delete("startAt");
                          setSearchParams(newParams.current);
                        }}
                      >
                        <ClearIcon />
                      </IconButton>
                    </div>
                  )}
                />
              </Grid>
              <Grid item xs={3}>
                <Controller
                  name="endAt"
                  control={control}
                  render={({ field }) => (
                    <div className="flex">
                      <DateTimePicker
                        slotProps={{
                          toolbar: {
                            toolbarFormat: "YYYY",
                            toolbarPlaceholder: "??",
                          },
                          actionBar: {
                            actions: ["clear", "accept"],
                          },
                        }}
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
                      <IconButton
                        onClick={() => {
                          setValue("endAt", null);
                          newParams.current.delete("endAt");
                          setSearchParams(newParams.current);
                        }}
                      >
                        <ClearIcon />
                      </IconButton>
                    </div>
                  )}
                />
              </Grid>

              {[
                {
                  type: "lpn",
                  label: "License Plate Number",
                },
                {
                  type: "status",
                  label: "Status",
                },
                {
                  type: "integratorName",
                  label: "Integrator",
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
              <TableCell>Plate Number</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>3rd party</TableCell>
              <TableCell>Amount (RM)</TableCell>
              <TableCell>Tax Info</TableCell>
              <TableCell>Error</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data?.data?.map((row) => (
              <TableRow key={row.businessTransactionId}>
                <TableCell>
                  <div className="w-[10rem]">
                    {moment(row.createdAt)
                      .local()
                      .format("DD/MM/yyyy hh:mm:ss A")}
                  </div>
                </TableCell>
                <TableCell>
                  <div className="w-[5rem]">{row.lpn}</div>
                </TableCell>
                <TableCell>
                  <Typography
                    className="w-[5rem]"
                    sx={{
                      color: row.status === "success" ? "#00afaa" : "#e4002b",
                    }}
                  >
                    {row.status}
                  </Typography>
                </TableCell>
                <TableCell>
                  <div className="w-[5rem]">{row.integratorName || "-"}</div>
                </TableCell>
                <TableCell>
                  <div className="w-[5rem]">{row.amount.toFixed(2) || "-"}</div>
                </TableCell>
                <TableCell>
                  <JsonViewer
                    className="w-[15rem]"
                    rootName={false}
                    displayDataTypes={false}
                    value={row.taxData}
                    defaultInspectControl={(_, __) => false}
                    collapseStringsAfterLength={false}
                  />
                </TableCell>
                <TableCell
                  style={{
                    width: "30px",
                    whiteSpace: "normal",
                    wordWrap: "break-word",
                  }}
                >
                  <div className="w-[30rem]">{row.error}</div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        <TablePagination
          rowsPerPageOptions={perPagesDefault.current}
          component="div"
          count={data?.metadata.totalData || 0}
          rowsPerPage={parseInt(
            currentQueryParameters.get("perPage") ||
              `${perPagesDefault.current[0]}`,
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

export default IntegratorTransactionsPage;
