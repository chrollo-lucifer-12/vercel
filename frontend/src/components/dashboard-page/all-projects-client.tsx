"use client";

import { useProject } from "@/hooks/use-project";
import ProjectCard from "../project-card";
import ScrollComponent from "../scroll-component";
import { useCallback, useEffect, useRef } from "react";

const AllProjectsClient = ({ name }: { name: string }) => {
  const rootRef = useRef(null);
  const elementRef = useRef<HTMLDivElement | null>(null);

  const { data, hasNextPage, fetchNextPage, isFetching } = useProject(
    name,
    16,
    { pages: [], pageParams: [] },
  );

  const observerCallback = useCallback(
    (entries: IntersectionObserverEntry[]) => {
      if (entries[0].isIntersecting && !isFetching) {
        fetchNextPage();
      }
    },
    [isFetching, fetchNextPage],
  );

  useEffect(() => {
    let observer: IntersectionObserver | null;
    if (rootRef?.current) {
      const options = {
        root: rootRef.current,
        rootMargin: "0px",
        scrollMargin: "50px",
        threshold: 1.0,
        delay: 2000,
      };
      observer = new IntersectionObserver(observerCallback, options);
      if (elementRef.current) {
        observer.observe(elementRef.current);
      }
    }
    return () => {
      if (observer) {
        observer.disconnect();
      }
    };
  }, [observerCallback]);

  const projects = data.pages.flatMap((page) => page ?? []);
  return (
    <div className="flex flex-col mt-10 gap-4">
      <p>Showing {projects?.length} projects</p>

      <div
        ref={rootRef}
        className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3
                      overflow-y-auto max-h-[600px] p-1"
      >
        {projects.map((project) => (
          <ProjectCard key={project.ID} project={project} />
        ))}

        <ScrollComponent
          hasNextPage={hasNextPage}
          isFetching={isFetching}
          sentinelRef={elementRef}
        />
      </div>
    </div>
  );
};

export default AllProjectsClient;
