"use client";

import { CaretDoubleDownIcon } from "@phosphor-icons/react";
import { Button } from "./ui/button";
import { AnimatePresence, motion } from "motion/react";

type ScrollComponentProps = {
  hasNextPage: boolean;
  isFetching: boolean;
};

const ScrollComponent = ({ hasNextPage, isFetching }: ScrollComponentProps) => {
  return (
    <div className="col-span-full flex justify-center py-6 absolute bottom-2 -translate-x-1/2 left-1/2 z-10">
      <AnimatePresence mode="wait">
        {hasNextPage ? (
          <motion.div
            key="scroll"
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            transition={{ duration: 0.25 }}
            className="flex flex-col items-center gap-2"
          >
            <div className="animate-[bounce_0.8s_infinite]">
              <Button variant="outline" size="icon" className="rounded-full">
                <CaretDoubleDownIcon />
              </Button>
            </div>

            <p className="text-xs text-muted-foreground">
              {isFetching
                ? "Loading more projects..."
                : "Scroll down to see more projects"}
            </p>
          </motion.div>
        ) : (
          <motion.p
            key="end"
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            transition={{ duration: 0.25 }}
            className="text-xs text-muted-foreground"
          >
            No more projects available
          </motion.p>
        )}
      </AnimatePresence>
    </div>
  );
};

export default ScrollComponent;
