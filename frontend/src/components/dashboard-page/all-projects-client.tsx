"use client";

import { useSearchProjects } from "@/hooks/use-project";
import ProjectCard from "../project-card";
import ScrollComponent from "../scroll-component";
import { useCallback, useEffect, useRef } from "react";
import ProjectEmpty from "../project-empty";
import { Skeleton } from "../ui/skeleton";

const SkeletonGrid = () => {
  return (
    <div
      className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3
              overflow-y-auto max-h-[600px] p-1"
    >
      {Array.from({ length: 10 }).map((_, i) => {
        return <Skeleton key={i} />;
      })}
    </div>
  );
};

const AllProjectsClient = ({ name }: { name: string }) => {
  const rootRef = useRef(null);
  const elementRef = useRef<HTMLDivElement | null>(null);

  const {
    data,
    hasNextPage,
    fetchNextPage,
    isFetching,
    isFetchingNextPage,
    isPending,
  } = useSearchProjects(name, 12, { pages: [], pageParams: [] });

  const observerCallback = useCallback(
    (entries: IntersectionObserverEntry[]) => {
      if (entries[0].isIntersecting && hasNextPage && !isFetchingNextPage) {
        fetchNextPage();
      }
    },
    [isFetchingNextPage, fetchNextPage, hasNextPage],
  );

  useEffect(() => {
    let observer: IntersectionObserver | null;
    if (rootRef?.current) {
      const options = {
        root: rootRef.current,
        rootMargin: "100px",
        threshold: 0.1,
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
  if (isPending) {
    return <SkeletonGrid />;
  }

  if (projects.length === 0) {
    return <ProjectEmpty />;
  }

  return (
    <div className="flex flex-col mt-10 gap-4">
      <p>Showing {projects?.length} projects</p>

      <div
        ref={rootRef}
        className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3
                      overflow-y-auto max-h-[550px] p-1"
      >
        {projects.map((project) => (
          <ProjectCard key={project.id} project={project} name={name} />
        ))}

        <ScrollComponent
          hasNextPage={hasNextPage}
          isFetching={isFetchingNextPage}
        />

        <div ref={elementRef} className="h-10 col-span-full" />
      </div>
    </div>
  );
};

export default AllProjectsClient;
