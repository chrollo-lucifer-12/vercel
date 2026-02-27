import { clientEnv } from "../env/client";
import { Deployment, DeploymentWithLogs } from "../types";
import { axiosInstance } from "./axios";

export const getDeployments = async (accessToken: string, slug: string) => {
  try {
    const res = await axiosInstance.get<Deployment[]>(
      `${clientEnv.NEXT_PUBLIC_GET_DEPLOYMENTS}/${slug}`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );
    console.log(res.data);
    return res.data;
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to get deployments");
  }
};

export const getDeployment = async (
  accessToken: string,
  deploymentId: string,
) => {
  try {
    const res = await axiosInstance.get<DeploymentWithLogs>(
      `${clientEnv.NEXT_PUBLIC_GET_DEPLOYMENT}/${deploymentId}`,
      {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      },
    );
    console.log(res.data);
    return res.data;
  } catch (err) {
    console.error(err);
    throw err instanceof Error ? err : new Error("Failed to get deployment");
  }
};
