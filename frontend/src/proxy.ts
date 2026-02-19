import { cookies } from "next/headers";
import { NextRequest, NextResponse } from "next/server";

function isRefreshExpired(expiry?: string) {
  if (!expiry) return true;

  return Date.now() > Number(expiry);
}

export async function proxy(req: NextRequest) {
  // const refreshExpiry = req.cookies.get("refresh_token_expires_at")?.value;

  // if (isRefreshExpired(refreshExpiry)) {
  //   console.log("deleting refresh token");
  //   const cookieStore = await cookies();
  //   cookieStore.delete("refresh_token");
  // }

  return NextResponse.next();
}
