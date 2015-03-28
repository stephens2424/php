<?php

/*
 *  * This file is part of the Symfony package.
 *   *
 *    * (c) Fabien Potencier <fabien@symfony.com>
 *     *
 *      * For the full copyright and license information, please view the LICENSE
 *       * file that was distributed with this source code.
 *
 *         LICENSE ->
 *          Copyright (c) 2004-2014 Fabien Potencier
 *
 *          Permission is hereby granted, free of charge, to any person obtaining a copy
 *          of this software and associated documentation files (the "Software"), to deal
 *          in the Software without restriction, including without limitation the rights
 *          to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *          copies of the Software, and to permit persons to whom the Software is furnished
 *          to do so, subject to the following conditions:
 *
 *          The above copyright notice and this permission notice shall be included in all
 *          copies or substantial portions of the Software.
 *
 *          THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *          IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *          FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *          AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *          LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *          OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 *          THE SOFTWARE.
 *
 *            */

/**
 *  * PropelLogger.
 *   *
 *    * @author Fabien Potencier <fabien.potencier@symfony-project.com>
 *     * @author William Durand <william.durand1@gmail.com>
 *      */
class PropelLogger implements \BasicLogger
{

  /**
   *      * {@inheritdoc}
   *           */
  public function log($message, $severity = null)
  {
    if (null !== $this->logger) {
      $message = is_string($message) ? $message : var_export($message, true);

      switch ($severity) {
      case 'alert':
        $this->logger->alert($message);
        break;
      case 'crit':
        $this->logger->critical($message);
        break;
      case 'err':
        $this->logger->error($message);
        break;
      case 'warning':
        $this->logger->warning($message);
        break;
      case 'notice':
        $this->logger->notice($message);
        break;
      case 'info':
        $this->logger->info($message);
        break;
      case 'debug':
      default:
        $this->logger->debug($message);
      }
    }
  }

  /**
   *      * Returns queries.
   *           *
   *                * @return array Queries
   *                     */
  public function getQueries()
  {
    return $this->queries;
  }
}


