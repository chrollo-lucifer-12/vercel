"use client";

import { FileSearchIcon, PlusIcon } from "@phosphor-icons/react";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "../ui/input-group";
import CreateProject from "../create-project-dialog";
import { Button } from "../ui/button";
import { useProject } from "@/hooks/use-project";
import { useState } from "react";
import AllProjects from "./all-projects";

const SearchBar = () => {
  const [name, setName] = useState("");
  const { data, isLoading } = useProject(name);

  return (
    <div className="flex flex-col gap-10 py-3">
      <div className="flex gap-2">
        <InputGroup className="max-w-full">
          <InputGroupInput placeholder="Search..." />
          <InputGroupAddon>
            <FileSearchIcon />
          </InputGroupAddon>
        </InputGroup>
        <CreateProject>
          <Button size={"icon"} className="rounded-full">
            <PlusIcon />
          </Button>
        </CreateProject>
      </div>
      <AllProjects isLoading={isLoading} projects={data?.pages[0].projects!} />
    </div>
  );
};

export default SearchBar;
