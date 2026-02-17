"use client";

import CustomInputGroup from "@/components/custom/input-group";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Field, FieldDescription, FieldGroup } from "@/components/ui/field";
import { useSignIn } from "@/hooks/use-auth";
import { EnvelopeIcon, PasswordIcon } from "@phosphor-icons/react";

const SigninPage = () => {
  const { mutate, isPending, data } = useSignIn();
  return (
    <Card>
      <CardHeader>
        <CardTitle>Login to Account</CardTitle>
      </CardHeader>
      <CardContent>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            const form = e.currentTarget;
            const formData = new FormData(form);
            mutate(formData);
          }}
        >
          <FieldGroup>
            <Field>
              <CustomInputGroup
                icon={EnvelopeIcon}
                id="email"
                placeholder="john@gmail.com"
                text="Email"
                type="email"
                required={true}
                error={data?.error?.email?.[0] ?? null}
              />
            </Field>
            <Field>
              <CustomInputGroup
                icon={PasswordIcon}
                id="password"
                placeholder=""
                text="Password"
                type="password"
                required={true}
                error={data?.error?.password?.[0] ?? null}
              />
            </Field>
            <FieldGroup>
              <Field>
                <Button type="submit" disabled={isPending}>
                  Sign In
                </Button>
                <FieldDescription className="px-6 text-center">
                  Don't have an account? <a href="/sigup">Sign up</a>
                </FieldDescription>
              </Field>
            </FieldGroup>
          </FieldGroup>
        </form>
      </CardContent>
    </Card>
  );
};

export default SigninPage;
