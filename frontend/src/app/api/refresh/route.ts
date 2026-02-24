import { refresh } from "@/lib/axios/auth";
import { RefreshTokenDetails } from "@/lib/types";
import { cookies } from "next/headers";
import { NextRequest, NextResponse } from "next/server";

export const GET = async (req: NextRequest) => {
  try {
    const cookieStore = await cookies();
    const refreshTokenCookie = cookieStore.get("refresh_token")?.value;

    if (!refreshTokenCookie) {
      return NextResponse.json(
        { success: false, error: "No refresh token found" },
        { status: 401 },
      );
    }

    const parseCookie = JSON.parse(refreshTokenCookie) as RefreshTokenDetails;
    const res = await refresh(parseCookie.refresh_token);
    return NextResponse.json(
      {
        success: true,
        access_token: res,
      },
      { status: 201 },
    );
  } catch (err) {
    console.error(err);
    return NextResponse.json(
      {
        success: false,
        error: err instanceof Error ? err.message : "Unknown error",
      },
      { status: 500 },
    );
  }
};
