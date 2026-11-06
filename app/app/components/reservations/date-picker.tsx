"use client";

import * as React from "react";
import { format } from "date-fns";
import { Calendar as CalendarIcon } from "lucide-react";

import { cn } from "~/lib/utils";
import { Button } from "~/components/ui/button";
import { Calendar } from "~/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "~/components/ui/popover";
import { Field, FieldLabel } from "../ui/field";

interface DatePickerSimpleProps {
  value?: Date;
  onChange?: (date: Date) => void;
}

export function DatePickerSimple({ value, onChange }: DatePickerSimpleProps) {
  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          id="date-picker-simple"
          className="justify-start font-normal"
        >
          {value ? format(value, "PPP") : <span>Pick a date</span>}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0" align="start">
        <Calendar
          mode="single"
          selected={value}
          onSelect={onChange}
          defaultMonth={value}
          disabled={(date) => date > new Date()}
          required
        />
      </PopoverContent>
    </Popover>
  );
}
