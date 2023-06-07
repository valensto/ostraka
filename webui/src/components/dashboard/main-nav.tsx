import { cn } from "@/lib/utils";
import ThemeToggle from "./theme-toggle";

export function MainNav({
  className,
  ...props
}: React.HTMLAttributes<HTMLElement>) {
  return (
    <nav className={cn("flex h-16 items-center px-4", className)} {...props}>
      <h1 className="font-extrabold">Ostraka UI</h1>
      <div className="ml-auto">
        <a
          href="https://doc.ostraka.io"
          className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
        >
          Documentation
        </a>
      </div>
      <div className="inline-block ml-4">
          <ThemeToggle  />
      </div>
    </nav>
  );
}
