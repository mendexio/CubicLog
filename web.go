// CubicLog Web UI v1.1.0 - Beautiful embedded dashboard with smart analytics
//
// ARCHITECTURE:
// This file contains the complete web interface embedded as a Go string constant.
// The dashboard is built with modern web technologies served from CDNs:
// - Alpine.js v3 for reactive behavior and state management
// - Tailwind CSS v4 for modern styling and responsive design
// - Font Awesome for comprehensive iconography
// - Inter font family for professional typography
//
// CORE FEATURES:
// - Real-time log viewing with auto-refresh (5-second intervals)
// - Advanced search and filtering capabilities with regex support
// - Dark/light mode toggle with localStorage persistence
// - Color-coded log categories using 22 Tailwind CSS colors
// - Responsive design optimized for mobile and desktop
// - CSV/JSON export functionality with filtered data
// - Comprehensive pagination with configurable page sizes
// - JSON syntax highlighting with collapsible structures
//
// SMART ANALYTICS DASHBOARD:
// - Real-time server health monitoring with color-coded status indicators
// - Smart alerts system that automatically appears for high error rates
// - Error rate trending with visual indicators (increasing/decreasing/stable)
// - Volume trend analysis with activity level indicators
// - Automated severity breakdown using smart pattern recognition
// - Log type distribution with expandable/collapsible interface
// - Live statistics with automatic refresh and data synchronization
//
// DESIGN PHILOSOPHY:
// The entire UI is self-contained with no build process required.
// CDN resources maintain the single-binary deployment philosophy while
// providing a professional, modern interface that rivals dedicated
// logging platforms. The dashboard emphasizes clarity, speed, and
// actionable insights over complex configuration.
package main

// webUI contains the complete HTML dashboard as an embedded string
const webUI = `<!DOCTYPE html>
<html lang="en" class="dark">
<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>CubicLog - A Modern Logging Dashboard</title>
    
    <!-- Alpine.js -->
    <script defer src="https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js"></script>
    
    <!-- Tailwind CSS -->
    <script src="https://cdn.tailwindcss.com"></script>
    
    <!-- Font Awesome -->
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css" rel="stylesheet" />
    
    <!-- Google Fonts -->
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap" rel="stylesheet" />
    
    <script>
      tailwind.config = {
        darkMode: "class",
        theme: {
          extend: {
            fontFamily: {
              sans: ["Inter", "system-ui", "sans-serif"],
            },
            colors: {
              background: "#0a0a0a",
              foreground: "#fafafa",
              card: "#111111",
              "card-foreground": "#fafafa",
              border: "#262626",
              input: "#171717",
              primary: "#3b82f6",
              "primary-foreground": "#f8fafc",
              secondary: "#1f1f1f",
              "secondary-foreground": "#fafafa",
              muted: "#171717",
              "muted-foreground": "#a3a3a3",
              accent: "#262626",
              "accent-foreground": "#fafafa",
              success: "#10b981",
              warning: "#f59e0b",
              error: "#ef4444",
              info: "#3b82f6",
            },
          },
        },
      };
    </script>
    
    <style>
      [x-cloak] { display: none !important; }
      body {
        font-family: "Inter", system-ui, sans-serif;
      }

      /* Light mode styles */
      .light {
        --bg-background: #ffffff;
        --bg-foreground: #0a0a0a;
        --bg-card: #f8fafc;
        --bg-card-foreground: #0a0a0a;
        --bg-border: #e2e8f0;
        --bg-input: #f1f5f9;
        --bg-muted: #f1f5f9;
        --bg-muted-foreground: #64748b;
        --bg-accent: #f1f5f9;
        --bg-accent-foreground: #0a0a0a;
        --bg-secondary: #f1f5f9;
        --bg-secondary-foreground: #0a0a0a;
      }

      .light body { background-color: var(--bg-background); color: var(--bg-foreground); }
      .light .bg-background { background-color: var(--bg-background); }
      .light .bg-card { background-color: var(--bg-card); }
      .light .border-border { border-color: var(--bg-border); }
      .light .bg-input { background-color: var(--bg-input); }
      .light .bg-muted { background-color: var(--bg-muted); }
      .light .text-muted-foreground { color: var(--bg-muted-foreground); }
      .light .hover\\:bg-accent:hover { background-color: var(--bg-accent); }
      .light .hover\\:text-foreground:hover { color: var(--bg-foreground) !important; }

      .dark body { background-color: #0a0a0a; color: #fafafa; }

      .sparkline { width: 60px; height: 20px; }
      .log-entry { transition: all 0.2s ease; }
      .log-entry:hover { background-color: #171717; }
      .light .log-entry:hover { background-color: #f8fafc; }

      .expandable-content { max-height: 0; overflow: hidden; transition: max-height 0.3s ease; }
      .expandable-content.expanded { max-height: 500px; }

      .status-indicator { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }
      .status-success { background-color: #10b981; }
      .status-warning { background-color: #f59e0b; }
      .status-error { background-color: #ef4444; }
      .status-info { background-color: #3b82f6; }

      .percentage-bar { height: 4px; border-radius: 2px; overflow: hidden; background-color: #262626; }
      .light .percentage-bar { background-color: #e2e8f0; }

      /* JSON syntax highlighting */
      .json-key { color: #60a5fa; }
      .json-string { color: #34d399; }
      .json-number { color: #fbbf24; }
      .json-boolean { color: #f87171; }
      .json-null { color: #9ca3af; }
      .json-punctuation { color: #d1d5db; }

      .light .json-key { color: #2563eb; }
      .light .json-string { color: #059669; }
      .light .json-number { color: #d97706; }
      .light .json-boolean { color: #dc2626; }
      .light .json-null { color: #6b7280; }
      .light .json-punctuation { color: #374151; }

      /* Date input styling for proper visibility in both themes */
      .date-input {
        color-scheme: dark;
        color: #fafafa;
      }
      
      .light .date-input {
        color-scheme: light;
        color: #0a0a0a;
      }

      /* Header button hover effects for proper visibility */
      .hover-button:hover {
        color: #fafafa; /* Light color for dark theme */
      }
      
      .light .hover-button:hover {
        color: #0a0a0a; /* Dark color for light theme */
      }
    </style>
</head>
<body class="bg-background text-foreground min-h-screen" x-data="cubiclogApp()" x-init="init()" x-cloak>
    <!-- Header -->
    <header class="border-b border-border bg-card">
        <div class="max-w-7xl mx-auto px-6 py-4">
            <div class="flex items-center justify-between">
                <div class="flex items-center space-x-4">
                    <div class="flex items-center space-x-2">
                        <i class="fas fa-cube text-primary text-xl"></i>
                        <h1 class="text-xl font-semibold">CubicLog</h1>
                    </div>
                </div>
                <div class="flex-1 flex justify-center">
                    <div class="text-sm text-muted-foreground font-mono" id="current-datetime">
                        <!-- Dynamic datetime will be inserted here -->
                    </div>
                </div>
                <div class="flex items-center space-x-4">
                    <button @click="manualRefresh()" 
                            :disabled="refreshing"
                            class="text-muted-foreground hover-button transition-colors disabled:opacity-50"
                            title="Refresh data">
                        <i class="fas fa-sync-alt" :class="refreshing ? 'animate-spin' : ''"></i>
                    </button>
                    <button @click="toggleTheme()" 
                            class="text-muted-foreground hover-button transition-colors"
                            title="Toggle theme">
                        <i class="fas fa-sun dark:hidden"></i>
                        <i class="fas fa-moon hidden dark:inline"></i>
                    </button>
                </div>
            </div>
        </div>
    </header>

    <div class="max-w-7xl mx-auto px-6 py-8">
        <!-- Analytics Section -->
        <div class="mb-8">
            <h2 class="text-2xl font-semibold mb-6">Analytics</h2>

            <!-- Smart Alerts (Only shown when there are alerts) -->
            <div class="bg-card border border-border rounded-lg mb-6" x-show="analytics.alerts.length > 0">
                <div class="px-6 py-4 border-b border-border">
                    <h3 class="text-lg font-semibold flex items-center">
                        <i class="fas fa-exclamation-triangle text-yellow-500 mr-2"></i>
                        Smart Alerts
                    </h3>
                </div>
                <div class="px-6 py-6">
                    <div class="space-y-4">
                        <template x-for="alert in analytics.alerts" :key="alert.type">
                            <div class="flex items-start space-x-3 p-4 rounded-lg border" 
                                 :class="alert.severity === 'high' ? 'bg-red-50 dark:bg-red-950/50 border-red-200 dark:border-red-800' : 
                                        alert.severity === 'medium' ? 'bg-yellow-50 dark:bg-yellow-950/50 border-yellow-200 dark:border-yellow-800' : 
                                        'bg-blue-50 dark:bg-blue-950/50 border-blue-200 dark:border-blue-800'">
                                <i class="fas fa-bell text-sm mt-1" 
                                   :class="alert.severity === 'high' ? 'text-red-500 dark:text-red-400' : 
                                          alert.severity === 'medium' ? 'text-yellow-500 dark:text-yellow-400' : 
                                          'text-blue-500 dark:text-blue-400'">
                                </i>
                                <div>
                                    <h4 class="text-sm font-semibold" x-text="alert.message"></h4>
                                    <p class="text-xs text-muted-foreground" x-text="alert.details"></p>
                                </div>
                            </div>
                        </template>
                    </div>
                </div>
            </div>

            <!-- Basic Metrics Row -->
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
                <!-- Total Logs Card -->
                <div class="bg-card border border-border rounded-lg p-6">
                    <div class="flex items-center justify-between mb-4">
                        <div>
                            <p class="text-muted-foreground text-sm">Total Logs</p>
                            <p class="text-2xl font-semibold" x-text="stats.total"></p>
                        </div>
                        <div class="text-success">
                            <i class="fas fa-database text-lg"></i>
                        </div>
                    </div>
                    <div class="flex items-center justify-between text-sm">
                        <span class="text-muted-foreground">All time</span>
                        <span class="text-green-600 font-medium">
                            <i class="fas fa-arrow-up text-xs"></i> Active
                        </span>
                    </div>
                </div>

                <!-- Last 24 Hours Card -->
                <div class="bg-card border border-border rounded-lg p-6">
                    <div class="flex items-center justify-between mb-4">
                        <div>
                            <p class="text-muted-foreground text-sm">Last 24 Hours</p>
                            <p class="text-2xl font-semibold" x-text="stats.recent"></p>
                        </div>
                        <div class="text-info">
                            <i class="fas fa-clock text-lg"></i>
                        </div>
                    </div>
                    <div class="flex items-center justify-between text-sm">
                        <span class="text-muted-foreground">Recent activity</span>
                        <span class="text-blue-600 font-medium">
                            <i class="fas fa-pulse text-xs"></i> Live
                        </span>
                    </div>
                </div>

                <!-- Volume Trend Card -->
                <div class="bg-card border border-border rounded-lg p-6">
                    <div class="flex items-center justify-between mb-4">
                        <div>
                            <p class="text-muted-foreground text-sm">Volume Trend</p>
                            <p class="text-2xl font-semibold capitalize" x-text="analytics.trends.volume_trend"></p>
                        </div>
                        <div class="w-5 h-5 rounded-full flex items-center justify-center" 
                             :class="analytics.trends.volume_trend === 'increasing' ? 'bg-blue-100 text-blue-800' : 
                                    analytics.trends.volume_trend === 'decreasing' ? 'bg-yellow-100 text-yellow-800' : 
                                    'bg-green-100 text-green-800'">
                            <i class="fas text-xs" :class="analytics.trends.volume_trend === 'increasing' ? 'fa-arrow-up' : 
                                                          analytics.trends.volume_trend === 'decreasing' ? 'fa-arrow-down' : 
                                                          'fa-equals'"></i>
                        </div>
                    </div>
                    <div class="flex items-center justify-between text-sm">
                        <span class="text-muted-foreground">Activity level</span>
                        <span class="font-medium" 
                              :class="analytics.trends.volume_trend === 'increasing' ? 'text-blue-600' : 
                                     analytics.trends.volume_trend === 'decreasing' ? 'text-yellow-600' : 
                                     'text-green-600'"
                              x-text="analytics.trends.volume_trend === 'increasing' ? 'High' : 
                                     analytics.trends.volume_trend === 'decreasing' ? 'Low' : 
                                     'Normal'">
                        </span>
                    </div>
                </div>

                <!-- Server Health Card -->
                <div class="bg-card border border-border rounded-lg p-6">
                    <div class="flex items-center justify-between mb-4">
                        <div>
                            <p class="text-muted-foreground text-sm">Server Health</p>
                            <p class="text-2xl font-semibold" 
                               :class="analytics.error_rate > 30 ? 'text-red-600' : 
                                      analytics.error_rate > 10 ? 'text-yellow-600' : 
                                      'text-green-600'"
                               x-text="analytics.error_rate > 30 ? 'Critical' : 
                                      analytics.error_rate > 10 ? 'Warning' : 
                                      'Healthy'">
                            </p>
                        </div>
                        <div class="w-4 h-4 rounded-full" 
                             :class="analytics.error_rate > 30 ? 'bg-red-500' : 
                                    analytics.error_rate > 10 ? 'bg-yellow-500' : 
                                    'bg-green-500'">
                        </div>
                    </div>
                    <div class="flex items-center justify-between text-sm">
                        <span class="text-muted-foreground">System status</span>
                        <span class="font-medium" 
                              :class="analytics.error_rate > 30 ? 'text-red-600' : 
                                     analytics.error_rate > 10 ? 'text-yellow-600' : 
                                     'text-green-600'"
                              x-text="analytics.error_rate > 30 ? 'Needs attention' : 
                                     analytics.error_rate > 10 ? 'Monitor closely' : 
                                     'All systems go'">
                        </span>
                    </div>
                </div>
            </div>



            <!-- Collapsible Log Distribution Card -->
            <div class="bg-card border border-border rounded-lg">
                <div class="px-6 py-4 border-b border-border">
                    <button @click="distributionExpanded = !distributionExpanded" 
                            class="flex items-center justify-between w-full text-left">
                        <div>
                            <h3 class="text-lg font-semibold">Log Type Distribution</h3>
                            <p class="text-muted-foreground text-sm">Breakdown by log type</p>
                        </div>
                        <i class="fas fa-chevron-down text-muted-foreground transform transition-transform duration-200" 
                           :class="distributionExpanded ? '' : '-rotate-90'"></i>
                    </button>
                </div>
                <div x-show="distributionExpanded" x-transition class="px-6 py-6">
                    <div class="space-y-4">
                        <template x-for="stat in dynamicStats" :key="stat.type">
                            <div class="flex items-center justify-between py-3 border-b border-border last:border-b-0">
                                <div class="flex items-center space-x-3">
                                    <span class="status-indicator" :style="'background-color: ' + stat.color"></span>
                                    <span class="text-sm font-medium" x-text="stat.label"></span>
                                </div>
                                <div class="flex items-center space-x-4">
                                    <div class="flex-1 w-24">
                                        <div class="percentage-bar">
                                            <div class="h-full" :style="'background-color: ' + stat.color + '; width: ' + Math.round((stat.count / stats.total) * 100) + '%'"></div>
                                        </div>
                                    </div>
                                    <div class="text-right min-w-0">
                                        <p class="text-sm font-semibold" x-text="stat.count"></p>
                                        <p class="text-xs text-muted-foreground" x-text="Math.round((stat.count / stats.total) * 100) + '%'" x-show="stats.total > 0"></p>
                                    </div>
                                </div>
                            </div>
                        </template>
                        <div x-show="dynamicStats.length === 0" class="text-center py-8 text-muted-foreground">
                            <i class="fas fa-inbox text-4xl mb-4 opacity-50"></i>
                            <p class="text-sm">No logs yet</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Search Section -->
        <div class="mb-8">
            <div class="bg-card border border-border rounded-lg p-6">
                <div class="flex flex-col lg:flex-row gap-4">
                    <div class="flex-1">
                        <div class="relative">
                            <i class="fas fa-search absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground"></i>
                            <input type="text"
                                   x-model="searchQuery"
                                   @input="applyFilters()"
                                   placeholder="Search logs..."
                                   class="w-full pl-10 pr-4 py-3 bg-input border border-border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent">
                        </div>
                    </div>
                    <div class="flex gap-3">
                        <select x-model="typeFilter"
                                @change="applyFilters()"
                                class="px-4 py-3 bg-input border border-border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary">
                            <option value="">All Levels</option>
                            <template x-for="type in uniqueTypes" :key="type">
                                <option :value="type" x-text="type.charAt(0).toUpperCase() + type.slice(1)"></option>
                            </template>
                        </select>
                        <input type="date"
                               x-model="selectedDate"
                               @change="applyFilters()"
                               class="px-4 py-3 bg-input border border-border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary date-input"
                               title="Filter by date">
                        <button @click="clearFilters()"
                                :disabled="clearing"
                                class="px-6 py-3 bg-primary text-primary-foreground rounded-lg hover:bg-primary/90 transition-colors disabled:opacity-50"
                                :class="clearing ? 'scale-95' : ''">
                            <i class="fas fa-times mr-2"></i>
                            Clear
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Log List Section -->
        <div class="mb-8">
            <!-- Loading -->
            <div x-show="loading" class="text-center py-12">
                <div class="flex items-center justify-center space-x-2">
                    <i class="fas fa-spinner animate-spin text-blue-500 text-xl"></i>
                    <p class="text-muted-foreground">Loading logs...</p>
                </div>
            </div>

            <!-- Logs -->
            <div x-show="!loading" class="bg-card border border-border rounded-lg overflow-hidden">
                <div class="border-b border-border px-6 py-4">
                    <h3 class="text-lg font-semibold">Recent Logs</h3>
                </div>

                <div class="divide-y divide-border">
                    <template x-for="log in filteredLogs" :key="log.id">
                        <div class="log-entry cursor-pointer" @click="toggleLogExpansion(log.id)">
                            <div class="px-6 py-4 flex items-center justify-between">
                                <div class="flex items-center space-x-4 flex-1">
                                    <span class="status-indicator" :style="'background-color: ' + getLogColor(log.header.color, log.header.type)"></span>
                                    <div class="flex-1">
                                        <div class="flex items-center space-x-3">
                                            <span class="text-sm font-mono text-muted-foreground" x-text="formatTime(log.timestamp)"></span>
                                            <span class="px-2 py-1 text-xs rounded-full" 
                                                  :class="getTypeBadgeClass(log.header.type, log.header.color)"
                                                  x-text="log.header.type.toUpperCase()"></span>
                                            <span class="text-sm text-muted-foreground" x-text="log.header.source" x-show="log.header.source"></span>
                                        </div>
                                        <p class="text-sm mt-1" x-text="log.header.title"></p>
                                        <p class="text-xs text-muted-foreground mt-1" x-text="log.header.description" x-show="log.header.description"></p>
                                    </div>
                                </div>
                                <i class="fas fa-chevron-down text-muted-foreground transform transition-transform duration-200"
                                   :class="expandedLogs.includes(log.id) ? 'rotate-180' : ''"></i>
                            </div>
                            <div x-show="expandedLogs.includes(log.id)" x-transition class="px-6 pb-4">
                                <div class="bg-muted rounded-lg p-4">
                                    <div x-show="log.body && Object.keys(log.body).length > 0">
                                        <pre class="text-xs overflow-x-auto" x-html="formatJSON(log.body)"></pre>
                                    </div>
                                    <div x-show="!log.body || Object.keys(log.body).length === 0" class="text-xs text-muted-foreground">
                                        No additional data
                                    </div>
                                </div>
                            </div>
                        </div>
                    </template>

                    <!-- Empty state -->
                    <div x-show="filteredLogs.length === 0 && !loading" class="text-center py-12">
                        <i class="fas fa-search text-4xl text-muted-foreground opacity-50 mb-4"></i>
                        <p class="text-muted-foreground">
                            <span x-show="!searchQuery && !typeFilter">Start sending logs to see them here</span>
                            <span x-show="searchQuery || typeFilter">No logs match your current filters</span>
                        </p>
                    </div>
                </div>
            </div>
        </div>

        <!-- Pagination -->
        <div x-show="!loading && totalPages > 1" class="flex items-center justify-between">
            <p class="text-sm text-muted-foreground">
                Showing <span class="font-medium" x-text="((currentPage - 1) * logsPerPage + 1)"></span> to
                <span class="font-medium" x-text="Math.min(currentPage * logsPerPage, totalLogs)"></span> of
                <span class="font-medium" x-text="totalLogs"></span> results
            </p>
            <div class="flex items-center space-x-2">
                <button @click="previousPage()"
                        :disabled="currentPage <= 1"
                        class="px-3 py-2 text-sm border border-border rounded-lg hover:bg-accent disabled:opacity-50 disabled:cursor-not-allowed">
                    <i class="fas fa-chevron-left mr-1"></i>
                    Previous
                </button>
                
                <!-- Logs per page dropdown -->
                <div class="flex items-center space-x-2">
                    <span class="text-sm text-muted-foreground">Show:</span>
                    <select x-model="logsPerPage" 
                            @change="changeLogsPerPage()"
                            class="px-3 py-2 text-sm border border-border rounded-lg bg-input hover:bg-accent focus:outline-none focus:ring-2 focus:ring-primary">
                        <option value="10">10</option>
                        <option value="25">25</option>
                        <option value="50">50</option>
                    </select>
                </div>
                
                <button @click="nextPage()"
                        :disabled="currentPage >= totalPages"
                        class="px-3 py-2 text-sm border border-border rounded-lg hover:bg-accent disabled:opacity-50 disabled:cursor-not-allowed">
                    Next
                    <i class="fas fa-chevron-right ml-1"></i>
                </button>
            </div>
        </div>
    </div>

    <!-- Footer -->
    <footer class="border-t border-border bg-card mt-16">
        <div class="max-w-7xl mx-auto px-6 py-6">
            <div class="text-center text-muted-foreground text-sm">
                Created with <i class="fas fa-heart text-red-500 mx-1"></i> by Mendex
                • <span id="current-year"></span>
            </div>
        </div>
    </footer>

    <script>
        function cubiclogApp() {
            return {
                // Data
                logs: [],
                filteredLogs: [],
                searchQuery: '',
                typeFilter: '',
                selectedDate: '',
                expandedLogs: [],
                loading: true,
                refreshing: false,
                clearing: false,
                stats: {
                    total: 0,
                    recent: 0,
                    monthly: 0
                },
                analytics: {
                    error_rate: 0,
                    severity_breakdown: {},
                    top_sources: [],
                    hourly_distribution: [],
                    alerts: [],
                    trends: {
                        error_trend: 'stable',
                        volume_trend: 'stable'
                    }
                },
                uniqueTypes: [],
                dynamicStats: [],
                // Pagination
                currentPage: 1,
                logsPerPage: 10,
                totalPages: 0,
                totalLogs: 0,
                // UI state
                distributionExpanded: false,

                async init() {
                    // Load logs per page preference from localStorage
                    const savedLogsPerPage = localStorage.getItem('cubiclog_logs_per_page');
                    if (savedLogsPerPage) {
                        this.logsPerPage = parseInt(savedLogsPerPage);
                    }
                    
                    await this.fetchLogs();
                    // Auto-refresh every 5 seconds
                    setInterval(() => this.fetchLogs(), 5000);
                },

                async fetchLogs() {
                    if (this.loading) {
                        // Initial load
                    }
                    
                    try {
                        // Always fetch all logs first to maintain uniqueTypes and get total count
                        let allLogsUrl = '/api/logs?limit=1000';
                        const allLogsResponse = await fetch(allLogsUrl);
                        const allLogs = await allLogsResponse.json();
                        this.logs = allLogs;
                        this.updateUniqueTypes();
                        
                        // Get total count for pagination
                        this.totalLogs = this.logs.length;
                        this.totalPages = Math.ceil(this.totalLogs / this.logsPerPage);
                        
                        // Build paginated URL for display
                        const offset = (this.currentPage - 1) * this.logsPerPage;
                        let url = '/api/logs?limit=' + this.logsPerPage + '&offset=' + offset;
                        if (this.searchQuery) url += '&q=' + encodeURIComponent(this.searchQuery);
                        if (this.typeFilter) url += '&type=' + encodeURIComponent(this.typeFilter);
                        if (this.selectedDate) url += '&from=' + this.selectedDate;
                        
                        const response = await fetch(url);
                        this.filteredLogs = await response.json();
                        this.updateStats();
                        await this.fetchAnalytics();
                        
                    } catch (error) {
                        console.error('Error fetching logs:', error);
                    } finally {
                        this.loading = false;
                    }
                },

                async fetchAnalytics() {
                    try {
                        const response = await fetch('/api/stats');
                        const data = await response.json();
                        
                        // Parse error rate from string percentage to number
                        const errorRate = parseFloat((data.error_rate_24h || '0%').replace('%', ''));
                        
                        // Map backend structure to frontend expectations
                        this.analytics = {
                            error_rate: errorRate,
                            severity_breakdown: data.severity_breakdown || {},
                            top_sources: data.top_sources ? data.top_sources.map(src => ({
                                source: src.name,
                                count: src.count
                            })) : [],
                            hourly_distribution: data.hourly_distribution || [],
                            alerts: Array.isArray(data.alerts) ? data.alerts.map(alert => ({
                                type: 'error_rate',
                                message: alert,
                                details: 'Automated detection based on recent log patterns',
                                severity: errorRate > 30 ? 'high' : errorRate > 15 ? 'medium' : 'low'
                            })) : [],
                            trends: {
                                error_trend: data.trends?.errors_increasing ? 'increasing' : 
                                           data.trends?.error_change < 0 ? 'decreasing' : 'stable',
                                volume_trend: data.trends?.spike_detected ? 'increasing' : 'stable'
                            }
                        };
                    } catch (error) {
                        console.error('Error fetching analytics:', error);
                    }
                },
                
                async manualRefresh() {
                    this.refreshing = true;
                    try {
                        await this.fetchLogs();
                        await new Promise(resolve => setTimeout(resolve, 500));
                    } catch (error) {
                        console.error('Error fetching logs:', error);
                    } finally {
                        this.refreshing = false;
                    }
                },

                applyFilters() {
                    this.currentPage = 1;
                    this.fetchLogs();
                },

                async clearFilters() {
                    this.clearing = true;
                    try {
                        this.searchQuery = '';
                        this.typeFilter = '';
                        this.selectedDate = '';
                        this.currentPage = 1;
                        await this.fetchLogs();
                        await new Promise(resolve => setTimeout(resolve, 300));
                    } catch (error) {
                        console.error('Error clearing filters:', error);
                    } finally {
                        this.clearing = false;
                    }
                },

                // Pagination methods
                goToPage(page) {
                    if (page >= 1 && page <= this.totalPages) {
                        this.currentPage = page;
                        this.fetchLogs();
                    }
                },

                previousPage() {
                    if (this.currentPage > 1) {
                        this.currentPage--;
                        this.fetchLogs();
                    }
                },

                nextPage() {
                    if (this.currentPage < this.totalPages) {
                        this.currentPage++;
                        this.fetchLogs();
                    }
                },
                changeLogsPerPage() {
                    // Save preference to localStorage
                    localStorage.setItem('cubiclog_logs_per_page', this.logsPerPage);
                    // Reset to first page and fetch logs
                    this.currentPage = 1;
                    this.fetchLogs();
                },

                // UI functions
                toggleTheme() {
                    const html = document.documentElement;
                    if (html.classList.contains('dark')) {
                        html.classList.remove('dark');
                        html.classList.add('light');
                        localStorage.setItem('theme', 'light');
                    } else {
                        html.classList.remove('light');
                        html.classList.add('dark');
                        localStorage.setItem('theme', 'dark');
                    }
                },

                toggleLogExpansion(logId) {
                    const index = this.expandedLogs.indexOf(logId);
                    if (index > -1) {
                        this.expandedLogs.splice(index, 1);
                    } else {
                        this.expandedLogs.push(logId);
                    }
                },

                updateUniqueTypes() {
                    const types = [...new Set(this.logs.map(log => log.header.type))];
                    this.uniqueTypes = types.sort();
                },

                updateStats() {
                    this.stats.total = this.logs.length;
                    
                    // Recent logs (last 24 hours)
                    const oneDayAgo = new Date(Date.now() - 24 * 60 * 60 * 1000);
                    this.stats.recent = this.logs.filter(log => 
                        new Date(log.timestamp) > oneDayAgo
                    ).length;
                    
                    // Monthly logs (last 30 days)
                    const oneMonthAgo = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000);
                    this.stats.monthly = this.logs.filter(log => 
                        new Date(log.timestamp) > oneMonthAgo
                    ).length;
                    
                    this.updateDynamicStats();
                },
                
                updateDynamicStats() {
                    // Count logs by type
                    const typeCounts = {};
                    this.logs.forEach(log => {
                        const type = log.header.type;
                        typeCounts[type] = (typeCounts[type] || 0) + 1;
                    });
                    
                    this.dynamicStats = [];
                    
                    // Create stats for all types using the colors from the logs themselves
                    for (const [type, count] of Object.entries(typeCounts)) {
                        const logOfThisType = this.logs.find(log => log.header.type === type);
                        const color = this.getHexColor(type, logOfThisType?.header.color);
                        
                        this.dynamicStats.push({
                            type: type,
                            count: count,
                            color: color,
                            label: type.charAt(0).toUpperCase() + type.slice(1)
                        });
                    }
                    
                    // Sort by count (descending)
                    this.dynamicStats.sort((a, b) => b.count - a.count);
                },

                getHexColor(type, color) {
                    const colorMap = {
                        'red': '#ef4444',
                        'green': '#10b981', 
                        'blue': '#3b82f6',
                        'yellow': '#f59e0b',
                        'purple': '#8b5cf6',
                        'pink': '#ec4899',
                        'indigo': '#6366f1',
                        'cyan': '#06b6d4',
                        'orange': '#f97316',
                        'emerald': '#10b981',
                        'lime': '#65a30d',
                        'teal': '#0d9488',
                        'sky': '#0ea5e9',
                        'violet': '#8b5cf6',
                        'fuchsia': '#d946ef',
                        'rose': '#f43f5e',
                        'slate': '#64748b'
                    };
                    
                    if (color && colorMap[color]) {
                        return colorMap[color];
                    }
                    
                    // Default based on type
                    switch (type) {
                        case 'error': return '#ef4444';
                        case 'warning': return '#f59e0b';
                        case 'info': return '#3b82f6';
                        case 'debug': return '#6b7280';
                        default: return '#64748b';
                    }
                },

                getStatusClass(type) {
                    switch (type) {
                        case 'error': return 'status-error';
                        case 'warning': return 'status-warning';
                        case 'info': return 'status-success';
                        case 'debug': return 'status-info';
                        default: return 'status-info';
                    }
                },

                getTypeBadgeClass(type, color) {
                    const baseClasses = 'transition-colors';
                    
                    if (color) {
                        return baseClasses + ' bg-' + color + '-100 text-' + color + '-800';
                    }
                    
                    switch (type) {
                        case 'error': return baseClasses + ' bg-error/10 text-error';
                        case 'warning': return baseClasses + ' bg-warning/10 text-warning';
                        case 'info': return baseClasses + ' bg-success/10 text-success';
                        case 'debug': return baseClasses + ' bg-info/10 text-info';
                        default: return baseClasses + ' bg-gray-100 text-gray-800';
                    }
                },
                getLogColor(color, type) {
                    // Simple mapping of Tailwind color names to CSS values
                    const colors = {
                        'red': '#ef4444', 'green': '#10b981', 'blue': '#3b82f6', 'yellow': '#f59e0b',
                        'orange': '#f97316', 'purple': '#a855f7', 'pink': '#ec4899', 'indigo': '#6366f1',
                        'cyan': '#06b6d4', 'gray': '#6b7280', 'slate': '#64748b', 'zinc': '#71717a',
                        'neutral': '#737373', 'stone': '#78716c', 'lime': '#65a30d', 'emerald': '#059669',
                        'teal': '#0d9488', 'sky': '#0ea5e9', 'violet': '#8b5cf6', 'fuchsia': '#d946ef',
                        'rose': '#f43f5e', 'gold': '#f59e0b'
                    };
                    
                    // Use provided color or default to slate
                    return colors[color] || colors['slate'];
                },

                formatTime(timestamp) {
                    return new Date(timestamp).toLocaleString();
                },

                formatJSON(obj) {
                    if (!obj) return '<span class="json-null">null</span>';
                    
                    const json = JSON.stringify(obj, null, 2);
                    return json
                        .replace(/(".*?"):/g, '<span class="json-key">$1</span>:')
                        .replace(/: (".*?")/g, ': <span class="json-string">$1</span>')
                        .replace(/: (\\d+)/g, ': <span class="json-number">$1</span>')
                        .replace(/: (true|false)/g, ': <span class="json-boolean">$1</span>')
                        .replace(/: (null)/g, ': <span class="json-null">$1</span>')
                        .replace(/([{}\\[\\],])/g, '<span class="json-punctuation">$1</span>');
                }
            }
        }

        // Initialize theme from localStorage
        const savedTheme = localStorage.getItem('theme') || 'dark';
        document.documentElement.classList.remove('light', 'dark');
        document.documentElement.classList.add(savedTheme);

        // Update datetime every second
        function updateDateTime() {
            const now = new Date();
            const options = {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit',
                hour12: false,
            };
            const datetimeElement = document.getElementById('current-datetime');
            if (datetimeElement) {
                datetimeElement.textContent = now.toLocaleString('en-US', options).replace(',', ' •');
            }
        }

        // Update year in footer
        document.addEventListener('DOMContentLoaded', function() {
            const yearElement = document.getElementById('current-year');
            if (yearElement) {
                yearElement.textContent = new Date().getFullYear();
            }
        });

        // Update datetime immediately and then every second
        updateDateTime();
        setInterval(updateDateTime, 1000);
    </script>
</body>
</html>`