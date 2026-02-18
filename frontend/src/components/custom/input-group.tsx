import { Icon, WarningCircleIcon } from "@phosphor-icons/react";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
  InputGroupText,
} from "../ui/input-group";
import { Badge } from "../ui/badge";

const CustomInputGroup = ({
  id,
  type,
  placeholder,
  text,
  icon: IconComponent,
  required = false,
  error,
}: {
  id: string;
  type: string;
  placeholder: string;
  text: string;
  icon: Icon;
  required?: boolean;
  error?: string | null;
}) => {
  return (
    <>
      <InputGroup
        className={`group

      rounded-lg
      transition-all duration-150
      ${error ? "group-hover:text-amber-600 group-focus-within:text-amber-600" : "group-hover:text-[#26b3a6] group-focus-within:text-[#26b3a6]"}
      shadow-md
      `}
      >
        <InputGroupInput
          id={id}
          type={type}
          placeholder={placeholder}
          required={required}
          name={id}
        />
        <InputGroupAddon align="block-start">
          <IconComponent
            weight="duotone"
            className={`${error ? "group-hover:text-amber-600 group-focus-within:text-amber-600" : "group-hover:text-[#26b3a6] group-focus-within:text-[#26b3a6]"} group-hover:hidden group-focus-within:hidden transition duration-150`}
          />
          <IconComponent
            weight="fill"
            className={`${error ? "group-hover:text-amber-600 group-focus-within:text-amber-600" : "group-hover:text-[#26b3a6] group-focus-within:text-[#26b3a6]"} hidden group-hover:block group-focus-within:block transition duration-150`}
          />
          <InputGroupText
            className={`${error ? "group-hover:text-amber-600 group-focus-within:text-amber-600" : "group-hover:text-[#26b3a6] group-focus-within:text-[#26b3a6]"} transition duration-150`}
          >
            {text}
          </InputGroupText>
        </InputGroupAddon>
      </InputGroup>
      {error && (
        <div className="w-fit">
          <Badge variant="destructive" className="text-amber-600">
            <WarningCircleIcon />
            {error}
          </Badge>
        </div>
      )}
    </>
  );
};

export default CustomInputGroup;
