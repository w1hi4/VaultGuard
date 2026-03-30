"use client";

import React, { useState, useEffect } from "react";
import { 
  Shield, 
  Search, 
  AlertTriangle, 
  CheckCircle, 
  FileCode, 
  Terminal, 
  Shield as Github, 
  RefreshCcw,
  ExternalLink,
  ChevronRight,
  ShieldCheck,
  Zap
} from "lucide-react";

type Finding = {
  RuleID: string;
  Description: string;
  Commit: string;
  File: string;
  LineNumber: number;
  Match: string;
  Severity: string;
};

export default function Dashboard() {
  const [repoPath, setRepoPath] = useState("");
  const [scanning, setScanning] = useState(false);
  const [findings, setFindings] = useState<Finding[]>([]);
  const [stats, setStats] = useState({ critical: 0, high: 0, medium: 0, low: 0 });

  const startScan = async () => {
    setScanning(true);
    setFindings([]);
    try {
      const res = await fetch(`/api/scan?path=${encodeURIComponent(repoPath)}`);
      const data = await res.json();
      if (data.findings) {
        setFindings(data.findings);
        processStats(data.findings);
      }
    } catch (err) {
      console.error("Scan failed", err);
    } finally {
      setScanning(false);
    }
  };

  const processStats = (data: Finding[]) => {
    const s = { critical: 0, high: 0, medium: 0, low: 0 };
    data.forEach(f => {
      const sev = f.Severity.toLowerCase() as keyof typeof s;
      if (s[sev] !== undefined) s[sev]++;
    });
    setStats(s);
  };

  return (
    <div className="min-h-screen p-8 bg-[#0a0a0c] text-white">
      {/* Navigation */}
      <nav className="flex justify-between items-center mb-12">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-blue-600 rounded-lg shadow-[0_0_15px_rgba(37,99,235,0.4)]">
            <Shield size={28} className="text-white" />
          </div>
          <h1 className="text-2xl font-bold tracking-tight">
            Vault<span className="text-blue-500">Guard</span>
          </h1>
        </div>
        <div className="flex gap-4">
           <button className="flex items-center gap-2 px-4 py-2 glass rounded-full text-sm font-medium hover:bg-white/10 transition-all">
             <Github size={18} />
             Open Source
           </button>
        </div>
      </nav>

      {/* Hero / Controls */}
      <div className="max-w-4xl mx-auto mb-16">
        <div className="text-center mb-10">
          <h2 className="text-4xl font-extrabold mb-4 bg-gradient-to-r from-white to-gray-500 bg-clip-text text-transparent">
            Secure Your Git History
          </h2>
          <p className="text-gray-400 text-lg">
            High-performance secret scanning with advanced entropy analysis.
          </p>
        </div>

        <div className="flex gap-4 p-2 glass rounded-2xl shadow-2xl">
          <div className="flex-1 flex items-center px-4 gap-3">
            <Search className="text-gray-500" size={20} />
            <input 
              type="text" 
              placeholder="Enter local repository path..." 
              className="bg-transparent w-full outline-none text-white py-3 font-medium placeholder:text-gray-600"
              value={repoPath}
              onChange={(e) => setRepoPath(e.target.value)}
            />
          </div>
          <button 
            onClick={startScan}
            disabled={scanning || !repoPath}
            className={`px-8 py-3 rounded-xl font-bold transition-all flex items-center gap-2 ${
              scanning ? 'bg-gray-700 animate-pulse' : 'bg-blue-600 hover:bg-blue-500 shadow-[0_0_20px_rgba(37,99,235,0.3)]'
            }`}
          >
            {scanning ? <RefreshCcw className="animate-spin" /> : <Zap size={18} />}
            {scanning ? "Scanning..." : "Run Analysis"}
          </button>
        </div>
      </div>

      {findings.length > 0 && (
        <div className="max-w-6xl mx-auto animate-in fade-in slide-in-from-bottom-5 duration-700">
          {/* Stats Grid */}
          <div className="grid grid-cols-4 gap-6 mb-10">
            {[
              { label: "Critical", count: stats.critical, color: "text-red-500", bg: "bg-red-500/10" },
              { label: "High", count: stats.high, color: "text-orange-500", bg: "bg-orange-500/10" },
              { label: "Medium", count: stats.medium, color: "text-yellow-500", bg: "bg-yellow-500/10" },
              { label: "Low", count: stats.low, color: "text-blue-500", bg: "bg-blue-500/10" },
            ].map((s, i) => (
              <div key={i} className={`p-6 rounded-2xl glass ${s.bg}`}>
                <p className={`text-sm font-bold uppercase tracking-wider mb-1 ${s.color}`}>{s.label}</p>
                <p className="text-4xl font-black">{s.count}</p>
              </div>
            ))}
          </div>

          {/* Results Table */}
          <div className="glass rounded-3xl overflow-hidden">
            <div className="p-6 border-b border-white/5 flex justify-between items-center bg-white/5">
              <h3 className="text-xl font-bold flex items-center gap-2">
                <AlertTriangle className="text-yellow-500" />
                Detection Breakdown
              </h3>
              <span className="px-3 py-1 bg-white/10 rounded-full text-xs font-mono text-gray-400">
                {findings.length} findings
              </span>
            </div>
            <div className="divide-y divide-white/5">
              {findings.map((f, i) => (
                <div key={i} className="p-6 hover:bg-white/[0.02] transition-colors group">
                  <div className="flex justify-between items-start mb-4">
                    <div className="flex gap-4 items-center">
                      <div className={`p-2 rounded-lg ${
                        f.Severity === "CRITICAL" ? "bg-red-500/20 text-red-500" : 
                        f.Severity === "HIGH" ? "bg-orange-500/20 text-orange-500" :
                        "bg-yellow-500/20 text-yellow-500"
                      }`}>
                         <ShieldCheck size={20} />
                      </div>
                      <div>
                        <h4 className="font-bold text-lg mb-0.5">{f.Description}</h4>
                        <p className="text-gray-500 text-sm font-mono flex items-center gap-2">
                          <FileCode size={14} />
                          {f.File}:{f.LineNumber}
                        </p>
                      </div>
                    </div>
                    <span className="px-3 py-1 glass rounded-lg text-[10px] font-bold tracking-widest uppercase opacity-60">
                      ID: {f.RuleID}
                    </span>
                  </div>
                  
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="p-4 bg-black/40 rounded-xl border border-white/5 group-hover:border-white/10 transition-colors">
                       <p className="text-[10px] uppercase font-bold text-gray-500 mb-2 flex items-center gap-1">
                         <Terminal size={10} />
                         Snippet Match
                       </p>
                       <code className="text-blue-400 font-mono text-sm break-all">
                         {f.Match}
                       </code>
                    </div>
                    <div className="p-4 bg-black/40 rounded-xl border border-white/5 group-hover:border-white/10 transition-colors">
                       <p className="text-[10px] uppercase font-bold text-gray-500 mb-2 flex items-center gap-1">
                         <Github size={10} />
                         Commit Hash
                       </p>
                       <code className="text-gray-400 font-mono text-xs">
                         {f.Commit}
                       </code>
                    </div>
                  </div>
                  
                  <div className="mt-4 flex justify-end">
                    <button className="text-xs text-blue-500 font-bold flex items-center gap-1 hover:underline">
                      View Remediation Guide
                      <ChevronRight size={14} />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Empty State / Success */}
      {!scanning && findings.length === 0 && (
        <div className="flex flex-col items-center justify-center mt-20 text-center">
          <div className="w-24 h-24 bg-green-500/10 rounded-full flex items-center justify-center mb-6 border border-green-500/20 shadow-[0_0_30px_rgba(34,197,94,0.1)]">
            <CheckCircle size={48} className="text-green-500" />
          </div>
          <h3 className="text-2xl font-bold mb-2">System Ready</h3>
          <p className="text-gray-500 max-w-sm">
            Enter a repository path above to begin deep analysis of your Git history.
          </p>
        </div>
      )}

      <footer className="mt-20 text-center pb-8 border-t border-white/5 pt-8 text-gray-600 text-sm">
         <p>© 2026 VaultGuard Open Source Project. Secure by default.</p>
      </footer>
    </div>
  );
}
