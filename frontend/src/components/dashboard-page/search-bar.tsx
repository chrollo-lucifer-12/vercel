import { FileSearchIcon } from "@phosphor-icons/react";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "../ui/input-group";

const SearchBar = () => {
  return (
    <InputGroup className="max-w-full">
      <InputGroupInput placeholder="Search..." />
      <InputGroupAddon>
        <FileSearchIcon />
      </InputGroupAddon>
    </InputGroup>
  );
};

export default SearchBar;
