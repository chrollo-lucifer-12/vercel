"use client";

import { FileSearchIcon, PlusIcon } from "@phosphor-icons/react";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "../ui/input-group";
import CreateProject from "../create-project-dialog";
import { Button } from "../ui/button";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useDebouncedCallback } from "use-debounce";

const SearchBar = () => {
  const searchParams = useSearchParams();
  const pathname = usePathname();
  const { replace } = useRouter();
  const handleSearch = useDebouncedCallback((name: string) => {
    const params = new URLSearchParams(searchParams);
    if (name) {
      params.set("name", name);
    } else {
      params.delete("name");
    }
    replace(`${pathname}?${params.toString()}`);
  }, 500);

  return (
    <div className="flex flex-col gap-10 py-3">
      <div className="flex gap-2">
        <InputGroup className="max-w-full">
          <InputGroupInput
            placeholder="Search..."
            onChange={(e) => {
              handleSearch(e.target.value);
            }}
            defaultValue={searchParams.get("name")?.toString()}
          />
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
    </div>
  );
};

export default SearchBar;
