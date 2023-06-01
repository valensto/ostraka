import { cn } from "@/lib/utils";

export function MainNav({
  className,
  ...props
}: React.HTMLAttributes<HTMLElement>) {
  return (
    <nav
      className={cn(
        "flex items-center justify-between space-x-4 lg:space-x-6",
        className
      )}
      {...props}
    >
      <h1 className="font-extrabold">Ostraka UI</h1>
      <a
        href="https://doc.ostraka.io"
        className="text-sm font-medium text-muted-foreground transition-colors hover:text-primary"
      >
        Documentation
      </a>
    </nav>
  );
}
