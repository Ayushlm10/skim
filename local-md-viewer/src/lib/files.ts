import { glob } from 'glob';
import { readFileSync, statSync, existsSync } from 'fs';
import { join, basename } from 'path';
import matter from 'gray-matter';

export interface MarkdownFile {
  path: string;
  relativePath: string;
  name: string;
  title: string;
  content: string;
  frontmatter: Record<string, unknown>;
  modifiedAt: Date;
  size: number;
}

export interface FileTreeNode {
  name: string;
  path: string;
  type: 'file' | 'directory';
  children?: FileTreeNode[];
}

// Get the docs directory from environment or use default
function getDocsDir(): string {
  const envDir = process.env.DOCS_DIR || import.meta.env.DOCS_DIR;
  if (envDir && existsSync(envDir)) {
    return envDir;
  }
  // Default to a 'docs' folder in current working directory
  const defaultDir = join(process.cwd(), 'docs');
  if (existsSync(defaultDir)) {
    return defaultDir;
  }
  // Fallback to current directory
  return process.cwd();
}

export async function getMarkdownFiles(): Promise<MarkdownFile[]> {
  const docsDir = getDocsDir();
  
  const files = await glob('**/*.md', {
    cwd: docsDir,
    ignore: ['node_modules/**', '.git/**', '**/node_modules/**'],
    absolute: false,
  });

  return files.map((relativePath) => {
    const fullPath = join(docsDir, relativePath);
    const raw = readFileSync(fullPath, 'utf-8');
    const { data, content } = matter(raw);
    const stats = statSync(fullPath);
    const name = basename(relativePath, '.md');
    
    return {
      path: fullPath,
      relativePath,
      name,
      title: (data.title as string) || formatTitle(name),
      content,
      frontmatter: data,
      modifiedAt: stats.mtime,
      size: stats.size,
    };
  }).sort((a, b) => a.relativePath.localeCompare(b.relativePath));
}

export async function getFileTree(): Promise<FileTreeNode[]> {
  const files = await getMarkdownFiles();
  const tree: FileTreeNode[] = [];
  
  for (const file of files) {
    const parts = file.relativePath.split('/');
    let current = tree;
    
    for (let i = 0; i < parts.length; i++) {
      const part = parts[i];
      const isFile = i === parts.length - 1;
      const path = parts.slice(0, i + 1).join('/');
      
      let node = current.find(n => n.name === part);
      
      if (!node) {
        node = {
          name: isFile ? file.title : part,
          path: isFile ? file.relativePath : path,
          type: isFile ? 'file' : 'directory',
          children: isFile ? undefined : [],
        };
        current.push(node);
      }
      
      if (!isFile && node.children) {
        current = node.children;
      }
    }
  }
  
  return sortTree(tree);
}

function sortTree(nodes: FileTreeNode[]): FileTreeNode[] {
  return nodes.sort((a, b) => {
    // Directories first
    if (a.type !== b.type) {
      return a.type === 'directory' ? -1 : 1;
    }
    return a.name.localeCompare(b.name);
  }).map(node => ({
    ...node,
    children: node.children ? sortTree(node.children) : undefined,
  }));
}

export async function getFileByPath(relativePath: string): Promise<MarkdownFile | null> {
  const files = await getMarkdownFiles();
  return files.find(f => f.relativePath === relativePath) || null;
}

function formatTitle(name: string): string {
  return name
    .replace(/[-_]/g, ' ')
    .replace(/\b\w/g, c => c.toUpperCase());
}

export function formatFileSize(bytes: number): string {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

export function formatDate(date: Date): string {
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
}
