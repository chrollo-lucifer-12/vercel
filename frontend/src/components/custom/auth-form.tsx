"use client";

import { Icon } from "@phosphor-icons/react";
import { Button } from "../ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../ui/card";
import { Field, FieldDescription, FieldGroup } from "../ui/field";
import CustomInputGroup from "./input-group";
import { Spinner } from "../ui/spinner";

type FieldConfig = {
  id: string;
  placeholder: string;
  text: string;
  type: string;
  icon: Icon;
  required?: boolean;
};

type AuthFormProps = {
  title: string;
  description?: string;
  fields: FieldConfig[];
  submitText: string;
  footerText?: React.ReactNode;
  onSubmit: (formData: FormData) => void;
  isPending?: boolean;
  errors?: Record<string, string[] | undefined>;
  loadingText: string;
};

const AuthForm = ({
  fields,
  onSubmit,
  submitText,
  title,
  description,
  errors,
  footerText,
  isPending,
  loadingText,
}: AuthFormProps) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        {description && <CardDescription>{description}</CardDescription>}
      </CardHeader>

      <CardContent>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            const formData = new FormData(e.currentTarget);
            onSubmit(formData);
          }}
        >
          <FieldGroup>
            {fields.map((field) => (
              <Field key={field.id}>
                <CustomInputGroup
                  {...field}
                  error={errors?.[field.id]?.[0] ?? null}
                />
              </Field>
            ))}

            <Field>
              <Button type="submit" disabled={isPending}>
                {isPending ? (
                  <>
                    <Spinner data-icon="inline-start" />
                    {loadingText}
                  </>
                ) : (
                  submitText
                )}
              </Button>

              {footerText && (
                <FieldDescription className="px-6 text-center ">
                  {footerText}
                  <a href="/verify">Verify Email</a>
                </FieldDescription>
              )}
            </Field>
          </FieldGroup>
        </form>
      </CardContent>
    </Card>
  );
};

export default AuthForm;
