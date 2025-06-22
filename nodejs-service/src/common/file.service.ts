import { Injectable } from '@nestjs/common';
import * as fs from 'fs';
import * as path from 'path';
import { Stream } from 'stream';

@Injectable()
export class FileService {
  private readonly publicStoragePath = path.resolve('./storage/public');
  private readonly privateStoragePath = path.resolve('./storage/private');

  constructor() {
    if (!fs.existsSync(this.publicStoragePath)) {
      fs.mkdirSync(this.publicStoragePath, { recursive: true });
    }
    if (!fs.existsSync(this.privateStoragePath)) {
      fs.mkdirSync(this.privateStoragePath, { recursive: true });
    }
  }

  async saveFile(
    fileStream: Stream,
    filename: string,
    userFolder: string,
  ): Promise<string> {
    const userFolderPath = path.join(this.publicStoragePath, userFolder);
    if (!fs.existsSync(userFolderPath)) {
      fs.mkdirSync(userFolderPath, { recursive: true });
    }

    const filePath = path.join(userFolderPath, filename);
    const writeStream = fs.createWriteStream(filePath);

    return new Promise((resolve, reject) => {
      fileStream.pipe(writeStream);
      fileStream.on('end', () => resolve(path.join(userFolder, filename)));
      fileStream.on('error', reject);
    });
  }

  async savePrivateFile(
    fileStream: Stream,
    filename: string,
    subfolder: string,
  ): Promise<string> {
    const folderPath = path.join(this.privateStoragePath, subfolder);
    if (!fs.existsSync(folderPath)) {
      fs.mkdirSync(folderPath, { recursive: true });
    }

    const filePath = path.join(folderPath, filename);
    const writeStream = fs.createWriteStream(filePath);

    return new Promise((resolve, reject) => {
      fileStream
        .pipe(writeStream)
        .on('finish', () => resolve(path.join(subfolder, filename)))
        .on('error', reject);
    });
  }

  getPrivateFileStream(relativePath: string): fs.ReadStream {
    const filePath = path.join(this.privateStoragePath, relativePath);
    if (!fs.existsSync(filePath)) {
      throw new Error('File not found');
    }
    return fs.createReadStream(filePath);
  }
} 