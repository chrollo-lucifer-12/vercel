"use client";

import { FileSearchIcon, PlusIcon } from "@phosphor-icons/react";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "../ui/input-group";
import CreateProject from "../create-project-dialog";
import { Button } from "../ui/button";

const SearchBar = () => {
  return (
    <div className="flex flex-row py-3 gap-2">
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
  );
};

export default SearchBar;
