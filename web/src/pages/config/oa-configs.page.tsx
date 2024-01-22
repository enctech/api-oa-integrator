import React, { useEffect, useRef, useState } from "react";
import {
  Button,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TablePagination,
  TableRow,
} from "@mui/material";
import { useMutation, useQuery } from "react-query";
import {
  deleteOAConfig,
  getOAConfigs,
  OAConfigResponse,
} from "../../api/config";
import { useNavigate } from "react-router-dom";
import { getOAHealth } from "../../api/health";
import { AdminOnly } from "../../components/auth-guard";
import DeleteIcon from "@mui/icons-material/Delete";
import IconButton from "@mui/material/IconButton";
import AlertDialog from "../../components/dialog";

const OAConfigsPage = () => {
  const navigate = useNavigate();

  const perPagesDefault = useRef([100, 500, 1000]);

  const [rowsPerPage, setRowsPerPage] = useState(perPagesDefault.current[0]);
  const [page, setPage] = useState(0);

  const { data, refetch } = useQuery("getOAConfigs", () => getOAConfigs());

  useEffect(() => {
    refetch().then();
  }, []);

  const handleChangePage = (_: any, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: { target: { value: string } }) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleRowClick = (id: string) => {
    navigate(`/oa-configs/${id}`);
  };

  const createNewConfig = () => {
    navigate(`/oa-configs/new`);
  };

  return (
    <div>
      <AdminOnly>
        <Button
          variant="contained"
          sx={{
            backgroundColor: "#3f3100",
            color: "#fff0bf",
            "&:hover": {
              backgroundColor: "#fff0bf",
              color: "#3f3100",
            },
          }}
          onClick={createNewConfig}
        >
          New
        </Button>
      </AdminOnly>
      <TableContainer component={Paper} className="mt-4">
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Facilities</TableCell>
              <TableCell>Devices</TableCell>
              <TableCell>Status</TableCell>
              <TableCell className="w-5 pl-1" />
            </TableRow>
          </TableHead>
          <TableBody>
            {data?.map((row) => (
              <OAConfigRow
                key={row.id}
                row={row}
                handleRowClick={handleRowClick}
                reFetch={refetch}
              />
            ))}
          </TableBody>
        </Table>
        <TablePagination
          rowsPerPageOptions={perPagesDefault.current}
          component="div"
          count={data?.length || 0}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={handleChangePage}
          onRowsPerPageChange={handleChangeRowsPerPage}
        />
      </TableContainer>
    </div>
  );
};

const OAConfigRow = ({
  row,
  handleRowClick,
  reFetch,
}: {
  row: OAConfigResponse;
  handleRowClick: (id: string) => void;
  reFetch: () => void;
}) => {
  const { data: oaHealth } = useQuery(
    "getOAHealth",
    () =>
      row.devices.length > 0 && row.facilities.length > 0
        ? getOAHealth({
            device: row.devices[0],
            facility: row.facilities[0],
          })
        : Promise.resolve(null),
    {},
  );

  const [showDeleteConfig, setShowDeleteConfigDialog] = useState(false);

  const { mutate: mutateDelete } = useMutation(
    "deleteOaConfig",
    deleteOAConfig,
    {
      onSuccess: () => reFetch(),
      onSettled: () => setShowDeleteConfigDialog(false),
    },
  );

  return (
    <>
      <TableRow
        className={"cursor-pointer"}
        key={row.id}
        onClick={() => handleRowClick(row.id)}
      >
        <TableCell>{row.name}</TableCell>
        <TableCell>{row.facilities?.join(", ")}</TableCell>
        <TableCell>{row.devices?.join(", ")}</TableCell>
        <TableCell>
          {oaHealth?.oa === "up" ? (
            <div
              className="w-5 h-5 rounded-full
                inline-flex items-center justify-center
                bg-green-500"
            ></div>
          ) : (
            <div
              className="w-5 h-5 rounded-full
                inline-flex items-center justify-center
                bg-red-500"
            ></div>
          )}
        </TableCell>
        <TableCell>
          <AdminOnly>
            <IconButton
              className="w-5 pl-1"
              onClick={(e) => {
                e.stopPropagation();
                setShowDeleteConfigDialog(true);
              }}
            >
              <DeleteIcon />
            </IconButton>
          </AdminOnly>
        </TableCell>
      </TableRow>
      <AlertDialog
        isOpen={showDeleteConfig}
        handleClose={() => setShowDeleteConfigDialog(false)}
        title={"Delete config"}
        description={"Are you sure you want to delete this configuration?"}
        buttons={[
          <Button
            key="cancel"
            onClick={() => setShowDeleteConfigDialog(false)}
            color="primary"
          >
            Cancel
          </Button>,
          <Button
            key="yes"
            onClick={() => {
              mutateDelete(row.id);
            }}
            color="primary"
            autoFocus
          >
            Delete
          </Button>,
        ]}
      />
    </>
  );
};

export default OAConfigsPage;
