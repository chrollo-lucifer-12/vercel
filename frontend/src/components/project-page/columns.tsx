"use client";

import { Deployment } from "@/lib/types";
import { ColumnDef } from "@tanstack/react-table";
import BuildDialog from "./build-dialog";
import { Button } from "../ui/button";
import { CalendarBlankIcon, DotsThreeCircleIcon } from "@phosphor-icons/react";
import { formatDate } from "@/lib/utils";
import { Badge } from "../ui/badge";
import { getQueryClient } from "@/lib/query-provider";

export const columns: ColumnDef<Deployment>[] = [
  {
    accessorKey: "sequence",
    header: "Sequence",
  },
  {
    accessorKey: "created_at",
    header: "Created At",
    cell: ({ row }) => {
      return (
        <Badge>
          <CalendarBlankIcon />
          {formatDate(row.original.created_at)}
        </Badge>
      );
    },
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => {
      const status = row.original.status;
      return (
        <Badge variant={`${status === "FAILED" ? "destructive" : "secondary"}`}>
          {status}
        </Badge>
      );
    },
  },
  {
    accessorKey: "action",
    header: "Logs",
    cell: ({ row }) => {
      const id = row.original.id;
      const queryClient = getQueryClient();
      return (
        <BuildDialog deploymentId={id}>
          <Button
            onClick={(e) => {
              queryClient.prefetchQuery({ queryKey: ["deployment", id] });
            }}
          >
            <DotsThreeCircleIcon />
          </Button>
        </BuildDialog>
      );
    },
  },
];
