import { useQueryState } from "nuqs";

export const useSearchQuery = (key: string) => {
  const [value, setValue] = useQueryState(key, {
    defaultValue: "",
    throttleMs: 2000,
    shallow: false,
  });

  return { value, setValue };
};
