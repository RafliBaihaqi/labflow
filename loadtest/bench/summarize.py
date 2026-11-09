"""Summarize a bench.sh run: median of each metric across rounds, per variant.

    python3 loadtest/bench/summarize.py loadtest/results/bench/<run-id>

Writes summary.json and summary.md next to the raw results.
"""

import json
import re
import statistics
import sys
from pathlib import Path

UNITS = {"B": 1, "KiB": 1024, "MiB": 1024**2, "GiB": 1024**3, "kB": 1000, "MB": 1000**2, "GB": 1000**3}


def mem_bytes(text):
    match = re.match(r"([\d.]+)\s*([A-Za-z]+)", text.strip())
    return float(match.group(1)) * UNITS[match.group(2)] if match else None


def stats_from_csv(path):
    cpu, mem = [], []
    for line in path.read_text().splitlines():
        parts = line.split(";")
        if len(parts) != 2 or not parts[0].endswith("%"):
            continue
        cpu.append(float(parts[0].rstrip("%")))
        used = mem_bytes(parts[1].split("/")[0])
        if used is not None:
            mem.append(used / 1024**2)
    if not cpu:
        return {}
    return {"cpu_avg_pct": statistics.mean(cpu), "cpu_max_pct": max(cpu),
            "mem_avg_mib": statistics.mean(mem), "mem_max_mib": max(mem)}


def k6_metrics(path):
    m = json.loads(path.read_text())["metrics"]

    def trend(scenario, stat):
        t = m.get(f"http_req_duration{{scenario:{scenario}}}")
        return t["values"][stat] if t else None

    return {
        "browse_p50_ms": trend("browse", "med"),
        "browse_p95_ms": trend("browse", "p(95)"),
        "browse_p99_ms": trend("browse", "p(99)"),
        "submit_p95_ms": trend("submit", "p(95)"),
        "login_p95_ms": trend("login", "p(95)"),
        "requests": m["http_reqs"]["values"]["count"],
        "rps": m["http_reqs"]["values"]["rate"],
        "failed_pct": m["http_req_failed"]["values"]["rate"] * 100,
        "reservations_created": m.get("reservations_created", {}).get("values", {}).get("count", 0),
    }


def main(run_dir):
    run_dir = Path(run_dir)
    meta = json.loads((run_dir / "meta.json").read_text())
    summary = {"meta": meta, "variants": {}}
    for label, info in meta["variants"].items():
        rounds = []
        for k6_file in sorted(run_dir.glob(f"{label}-r*.json")):
            if not re.fullmatch(rf"{re.escape(label)}-r\d+\.json", k6_file.name):
                continue
            row = k6_metrics(k6_file)
            row.update(stats_from_csv(k6_file.with_name(k6_file.stem + "-stats.csv")))
            rounds.append(row)
        medians = {}
        for key in rounds[0] if rounds else []:
            values = [r[key] for r in rounds if r.get(key) is not None]
            medians[key] = statistics.median(values) if values else None
        summary["variants"][label] = {**info, "rounds": rounds, "median": medians}

    (run_dir / "summary.json").write_text(json.dumps(summary, indent=2))

    rows = [
        ("Browse p50 (ms)", "browse_p50_ms", 1), ("Browse p95 (ms)", "browse_p95_ms", 1),
        ("Browse p99 (ms)", "browse_p99_ms", 1), ("Submit p95 (ms)", "submit_p95_ms", 1),
        ("Login p95 (ms)", "login_p95_ms", 0), ("Throughput (req/s)", "rps", 0),
        ("Failed requests (%)", "failed_pct", 2), ("API CPU avg (%)", "cpu_avg_pct", 0),
        ("API memory avg (MiB)", "mem_avg_mib", 0), ("API memory max (MiB)", "mem_max_mib", 0),
    ]
    labels = list(summary["variants"])
    warmup = meta.get("warmup_rounds", 0)
    lines = [f"Median of {meta['rounds']} rounds"
             + (f" (after {warmup} unrecorded warm-up round{'s' if warmup != 1 else ''})" if warmup else "")
             + f" at RATE={meta['rate']}, DURATION={meta['duration']}. "
             f"Host: {meta['host']}. API limited to {meta['api_cpus']} CPUs / {meta['api_memory']}.", "",
             "| Metric | " + " | ".join(labels) + " |", "| --- |" + " ---: |" * len(labels)]
    for title, key, digits in rows:
        cells = []
        for label in labels:
            value = summary["variants"][label]["median"].get(key)
            cells.append("-" if value is None else f"{value:,.{digits}f}")
        lines.append(f"| {title} | " + " | ".join(cells) + " |")
    for title, key in (("Binary size (MB)", "binary_bytes"), ("Image size (MB)", "image_bytes")):
        lines.append(f"| {title} | " + " | ".join(
            f"{summary['variants'][label][key] / 1e6:.1f}" for label in labels) + " |")
    (run_dir / "summary.md").write_text("\n".join(lines) + "\n")
    print("\n".join(lines))


if __name__ == "__main__":
    main(sys.argv[1])
